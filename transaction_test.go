package web3

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestGeneratePrivateKey(t *testing.T) {
	key, err := GeneratePrivateKey()
	if err != nil {
		t.Fatalf("GeneratePrivateKey() error = %v", err)
	}
	if key == nil {
		t.Fatal("GeneratePrivateKey() returned nil")
	}
}

func TestPrivateKeyRoundTrip(t *testing.T) {
	key, err := GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	hexKey := PrivateKeyToHex(key)
	if !strings.HasPrefix(hexKey, "0x") {
		t.Errorf("PrivateKeyToHex() = %q, want 0x prefix", hexKey)
	}
	recovered, err := PrivateKeyFromHex(hexKey)
	if err != nil {
		t.Fatalf("PrivateKeyFromHex() error = %v", err)
	}
	if PrivateKeyToHex(recovered) != hexKey {
		t.Error("round-trip private key mismatch")
	}
}

func TestPrivateKeyFromHexWithout0x(t *testing.T) {
	key, _ := GeneratePrivateKey()
	hexKey := PrivateKeyToHex(key)
	recovered, err := PrivateKeyFromHex(hexKey[2:]) // strip 0x
	if err != nil {
		t.Fatalf("PrivateKeyFromHex() without 0x error = %v", err)
	}
	if PrivateKeyToHex(recovered) != hexKey {
		t.Error("recovered key mismatch when 0x stripped")
	}
}

func TestPrivateKeyFromHexInvalid(t *testing.T) {
	_, err := PrivateKeyFromHex("0xinvalidhex")
	if err == nil {
		t.Error("expected error for invalid hex key")
	}
}

func TestPrivateKeyFromHexInvalidKey(t *testing.T) {
	// Valid hex but not a valid EC private key (all zeros)
	_, err := PrivateKeyFromHex("0x0000000000000000000000000000000000000000000000000000000000000000")
	if err == nil {
		t.Error("expected error for all-zero private key")
	}
}

func TestPrivateKeyToAddress(t *testing.T) {
	key, err := GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	addr := PrivateKeyToAddress(key)
	if !IsAddress(addr) {
		t.Errorf("PrivateKeyToAddress() = %q, not a valid address", addr)
	}
}

func TestPrivateKeyToAddressDeterministic(t *testing.T) {
	const knownHex = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	key, err := PrivateKeyFromHex(knownHex)
	if err != nil {
		t.Fatal(err)
	}
	addr := PrivateKeyToAddress(key)
	const want = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	if !strings.EqualFold(addr, want) {
		t.Errorf("PrivateKeyToAddress() = %q, want %q", addr, want)
	}
}

func TestNewTransactionParams(t *testing.T) {
	tp := NewTransactionParams()
	if tp == nil {
		t.Fatal("NewTransactionParams() returned nil")
	}
	if tp.Value == nil || tp.Value.Cmp(big.NewInt(0)) != 0 {
		t.Error("default Value should be 0")
	}
	if tp.ChainID == nil || tp.ChainID.Int64() != 1 {
		t.Error("default ChainID should be mainnet (1)")
	}
}

func TestTransactionParamsChaining(t *testing.T) {
	to := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	gasPrice := big.NewInt(20000000000)

	tp := NewTransactionParams().
		SetTo(to).
		SetGas(21000).
		SetGasPrice(gasPrice).
		SetNonce(5).
		SetChainID(ChainMainnet)

	if tp.To != to {
		t.Errorf("To = %q, want %q", tp.To, to)
	}
	if tp.Gas != 21000 {
		t.Errorf("Gas = %d, want 21000", tp.Gas)
	}
	if tp.GasPrice.Cmp(gasPrice) != 0 {
		t.Errorf("GasPrice = %v, want %v", tp.GasPrice, gasPrice)
	}
	if tp.Nonce != 5 {
		t.Errorf("Nonce = %d, want 5", tp.Nonce)
	}
}

func TestSetValueInEther(t *testing.T) {
	tp := NewTransactionParams().SetValueInEther("1")
	expected, _ := ToWei("1", Ether)
	if tp.Value.Cmp(expected) != 0 {
		t.Errorf("SetValueInEther(1) = %v, want %v", tp.Value, expected)
	}
}

func TestSetValueInWei(t *testing.T) {
	tp := NewTransactionParams().SetValueInWei("1000000000000000000")
	expected, _ := new(big.Int).SetString("1000000000000000000", 10)
	if tp.Value.Cmp(expected) != 0 {
		t.Errorf("SetValueInWei() = %v, want %v", tp.Value, expected)
	}
}

func TestSetGasPriceInGwei(t *testing.T) {
	tp := NewTransactionParams().SetGasPriceInGwei("20")
	expected, _ := ToWei("20", Gwei)
	if tp.GasPrice.Cmp(expected) != 0 {
		t.Errorf("SetGasPriceInGwei(20) = %v, want %v", tp.GasPrice, expected)
	}
}

func TestSetDataFromHex(t *testing.T) {
	tp := NewTransactionParams().SetDataFromHex("0xdeadbeef")
	if len(tp.Data) != 4 {
		t.Errorf("SetDataFromHex: data len = %d, want 4", len(tp.Data))
	}
	if tp.Data[0] != 0xde || tp.Data[1] != 0xad {
		t.Errorf("SetDataFromHex: unexpected data %x", tp.Data)
	}
}

func TestNewEIP1559TransactionParams(t *testing.T) {
	tp := NewEIP1559TransactionParams()
	if tp == nil {
		t.Fatal("NewEIP1559TransactionParams() returned nil")
	}
	if tp.ChainID == nil || tp.ChainID.Int64() != 1 {
		t.Error("default ChainID should be mainnet (1)")
	}
}

func TestSignTransaction(t *testing.T) {
	key, err := GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	tp := NewTransactionParams().
		SetTo("0x742d35Cc6634C0532925a3b844Bc454e4438f44e").
		SetGas(21000).
		SetGasPrice(big.NewInt(20000000000)).
		SetValue(big.NewInt(0)).
		SetNonce(0).
		SetChainID(ChainMainnet)

	signed, err := SignTransaction(tp, key)
	if err != nil {
		t.Fatalf("SignTransaction() error = %v", err)
	}
	if !strings.HasPrefix(signed.Hash, "0x") {
		t.Errorf("signed Hash = %q, want 0x prefix", signed.Hash)
	}
	if !strings.HasPrefix(signed.Raw, "0x") {
		t.Errorf("signed Raw = %q, want 0x prefix", signed.Raw)
	}
}

func TestSignTransactionDeployment(t *testing.T) {
	// Empty To is now allowed — produces a contract deployment tx.
	key, _ := GeneratePrivateKey()
	tp := NewTransactionParams().
		SetGas(500000).
		SetGasPrice(big.NewInt(20000000000)).
		SetData([]byte{0x60, 0x80, 0x60, 0x40}).
		SetNonce(0).
		SetChainID(ChainMainnet)

	signed, err := SignTransaction(tp, key)
	if err != nil {
		t.Fatalf("SignTransaction() deployment error = %v", err)
	}
	if !strings.HasPrefix(signed.Hash, "0x") {
		t.Errorf("deployment Hash = %q, want 0x prefix", signed.Hash)
	}
}

func TestSignTransactionMissingGasPrice(t *testing.T) {
	key, _ := GeneratePrivateKey()
	tp := NewTransactionParams().
		SetTo("0x742d35Cc6634C0532925a3b844Bc454e4438f44e").
		SetGas(21000)
	_, err := SignTransaction(tp, key)
	if err == nil {
		t.Error("expected error when GasPrice is missing")
	}
}

func TestSignTransactionMissingGas(t *testing.T) {
	key, _ := GeneratePrivateKey()
	tp := NewTransactionParams().
		SetTo("0x742d35Cc6634C0532925a3b844Bc454e4438f44e").
		SetGasPrice(big.NewInt(1))
	_, err := SignTransaction(tp, key)
	if err == nil {
		t.Error("expected error when Gas is zero")
	}
}

func TestSignEIP1559Transaction(t *testing.T) {
	key, err := GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	tp := NewEIP1559TransactionParams()
	tp.To = "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	tp.Gas = 21000
	tp.MaxFeePerGas = big.NewInt(30000000000)
	tp.MaxPriorityFeePerGas = big.NewInt(2000000000)
	tp.Value = big.NewInt(0)

	signed, err := SignEIP1559Transaction(tp, key)
	if err != nil {
		t.Fatalf("SignEIP1559Transaction() error = %v", err)
	}
	if !strings.HasPrefix(signed.Hash, "0x") {
		t.Errorf("Hash = %q, want 0x prefix", signed.Hash)
	}
}

func TestSignEIP1559MissingTo(t *testing.T) {
	key, _ := GeneratePrivateKey()
	tp := NewEIP1559TransactionParams()
	tp.Gas = 21000
	tp.MaxFeePerGas = big.NewInt(1)
	tp.MaxPriorityFeePerGas = big.NewInt(1)
	_, err := SignEIP1559Transaction(tp, key)
	if err == nil {
		t.Error("expected error when To is missing")
	}
}

func TestSignEIP1559MissingMaxFee(t *testing.T) {
	key, _ := GeneratePrivateKey()
	tp := NewEIP1559TransactionParams()
	tp.To = "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	tp.Gas = 21000
	tp.MaxPriorityFeePerGas = big.NewInt(1)
	_, err := SignEIP1559Transaction(tp, key)
	if err == nil {
		t.Error("expected error when MaxFeePerGas is missing")
	}
}

func TestSignEIP1559MissingMaxPriorityFee(t *testing.T) {
	key, _ := GeneratePrivateKey()
	tp := NewEIP1559TransactionParams()
	tp.To = "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	tp.Gas = 21000
	tp.MaxFeePerGas = big.NewInt(1)
	_, err := SignEIP1559Transaction(tp, key)
	if err == nil {
		t.Error("expected error when MaxPriorityFeePerGas is missing")
	}
}

func TestSignEIP1559MissingGas(t *testing.T) {
	key, _ := GeneratePrivateKey()
	tp := NewEIP1559TransactionParams()
	tp.To = "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	tp.MaxFeePerGas = big.NewInt(1)
	tp.MaxPriorityFeePerGas = big.NewInt(1)
	_, err := SignEIP1559Transaction(tp, key)
	if err == nil {
		t.Error("expected error when Gas is zero")
	}
}

func TestRecoverSigner(t *testing.T) {
	key, err := GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	expectedAddr := PrivateKeyToAddress(key)

	tp := NewTransactionParams().
		SetTo("0x742d35Cc6634C0532925a3b844Bc454e4438f44e").
		SetGas(21000).
		SetGasPrice(big.NewInt(20000000000)).
		SetValue(big.NewInt(0)).
		SetNonce(0).
		SetChainID(ChainMainnet)

	signed, err := SignTransaction(tp, key)
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := RecoverSigner(signed.Raw)
	if err != nil {
		t.Fatalf("RecoverSigner() error = %v", err)
	}
	if !strings.EqualFold(recovered, expectedAddr) {
		t.Errorf("RecoverSigner() = %q, want %q", recovered, expectedAddr)
	}
}

func TestRecoverSignerWithout0x(t *testing.T) {
	key, _ := GeneratePrivateKey()
	tp := NewTransactionParams().
		SetTo("0x742d35Cc6634C0532925a3b844Bc454e4438f44e").
		SetGas(21000).
		SetGasPrice(big.NewInt(20000000000)).
		SetNonce(0).
		SetChainID(ChainMainnet)
	signed, _ := SignTransaction(tp, key)

	// Strip the 0x prefix before passing to RecoverSigner.
	recovered, err := RecoverSigner(signed.Raw[2:])
	if err != nil {
		t.Fatalf("RecoverSigner() without 0x error = %v", err)
	}
	expected := PrivateKeyToAddress(key)
	if !strings.EqualFold(recovered, expected) {
		t.Errorf("RecoverSigner() = %q, want %q", recovered, expected)
	}
}

func TestRecoverSignerInvalidHex(t *testing.T) {
	_, err := RecoverSigner("0xnotvalidhex")
	if err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestRecoverSignerHomestead(t *testing.T) {
	// A Homestead (pre-EIP-155) transaction has ChainID=0 in the decoded tx.
	key, err := GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	ethTx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    big.NewInt(0),
		Gas:      21000,
		GasPrice: big.NewInt(20000000000),
	})
	signer := types.HomesteadSigner{}
	signed, err := types.SignTx(ethTx, signer, key)
	if err != nil {
		t.Fatal(err)
	}
	rawBytes, err := rlp.EncodeToBytes(signed)
	if err != nil {
		t.Fatal(err)
	}
	rawHex := fmt.Sprintf("0x%x", rawBytes)

	recovered, err := RecoverSigner(rawHex)
	if err != nil {
		t.Fatalf("RecoverSigner(Homestead) error = %v", err)
	}
	expected := PrivateKeyToAddress(key)
	if !strings.EqualFold(recovered, expected) {
		t.Errorf("RecoverSigner(Homestead) = %q, want %q", recovered, expected)
	}
}

func TestRecoverSignerInvalidRLP(t *testing.T) {
	// Valid hex but not a valid RLP-encoded transaction.
	_, err := RecoverSigner("0xdeadbeef")
	if err == nil {
		t.Error("expected error for invalid RLP bytes")
	}
}

func TestCreateContractDeployment(t *testing.T) {
	key, _ := GeneratePrivateKey()
	params := NewTransactionParams().
		SetGas(500000).
		SetGasPrice(big.NewInt(20000000000)).
		SetNonce(0).
		SetChainID(ChainMainnet)

	bytecode := []byte{0x60, 0x80, 0x60, 0x40}
	constructor := []byte{0x00, 0x01}

	signed, err := CreateContractDeployment(bytecode, constructor, key, params)
	if err != nil {
		t.Fatalf("CreateContractDeployment() error = %v", err)
	}
	if !strings.HasPrefix(signed.Hash, "0x") {
		t.Errorf("Hash = %q, want 0x prefix", signed.Hash)
	}
	// To should be empty for contract deployment
	if params.To != "" {
		t.Errorf("params.To = %q after deployment, want empty", params.To)
	}
}

func TestCreateContractDeploymentNoConstructor(t *testing.T) {
	key, _ := GeneratePrivateKey()
	params := NewTransactionParams().
		SetGas(500000).
		SetGasPrice(big.NewInt(20000000000)).
		SetNonce(0).
		SetChainID(ChainMainnet)

	signed, err := CreateContractDeployment([]byte{0x60, 0x80}, nil, key, params)
	if err != nil {
		t.Fatalf("CreateContractDeployment() without constructor error = %v", err)
	}
	if !strings.HasPrefix(signed.Hash, "0x") {
		t.Errorf("Hash = %q, want 0x prefix", signed.Hash)
	}
}

func TestCreateContractCall(t *testing.T) {
	key, _ := GeneratePrivateKey()
	params := NewTransactionParams().
		SetGas(100000).
		SetGasPrice(big.NewInt(20000000000)).
		SetNonce(0).
		SetChainID(ChainMainnet)

	signed, err := CreateContractCall("0x742d35Cc6634C0532925a3b844Bc454e4438f44e", []byte{0xde, 0xad}, key, params)
	if err != nil {
		t.Fatalf("CreateContractCall() error = %v", err)
	}
	if !strings.HasPrefix(signed.Hash, "0x") {
		t.Errorf("Hash = %q, want 0x prefix", signed.Hash)
	}
}

func TestCreateContractCallEmptyAddress(t *testing.T) {
	key, _ := GeneratePrivateKey()
	params := NewTransactionParams().
		SetGas(100000).
		SetGasPrice(big.NewInt(1))
	_, err := CreateContractCall("", []byte{0x01}, key, params)
	if err == nil {
		t.Error("expected error when contract address is empty")
	}
}

func TestEncodeABIAddress(t *testing.T) {
	data, err := EncodeABI("transfer(address,uint256)",
		"0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		big.NewInt(1000),
	)
	if err != nil {
		t.Fatalf("EncodeABI() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeABI() returned empty data")
	}
}

func TestEncodeABIBigInt(t *testing.T) {
	data, err := EncodeABI("setValue(uint256)", big.NewInt(42))
	if err != nil {
		t.Fatalf("EncodeABI() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeABI() returned empty data")
	}
}

func TestEncodeABIUint64(t *testing.T) {
	data, err := EncodeABI("setUint(uint64)", uint64(100))
	if err != nil {
		t.Fatalf("EncodeABI() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeABI() returned empty data")
	}
}

func TestEncodeABIBytes(t *testing.T) {
	data, err := EncodeABI("setBytes(bytes)", []byte{0x01, 0x02, 0x03})
	if err != nil {
		t.Fatalf("EncodeABI() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeABI() returned empty data")
	}
}

func TestEncodeABIBool(t *testing.T) {
	data, err := EncodeABI("setFlag(bool)", true)
	if err != nil {
		t.Fatalf("EncodeABI() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeABI() returned empty data")
	}
}

func TestEncodeABINonAddressString(t *testing.T) {
	// A non-address string gets typed as "string" in the ABI
	_, err := EncodeABI("setName(string)", "hello")
	// Result depends on library; just ensure no panic
	_ = err
}

func TestEncodeABIUnsupportedType(t *testing.T) {
	_, err := EncodeABI("setUint8(uint8)", uint8(42))
	if err == nil {
		t.Error("expected error for unsupported param type uint8")
	}
}

func TestRandomNonce(t *testing.T) {
	n1 := RandomNonce()
	n2 := RandomNonce()
	_ = n1
	_ = n2
}
