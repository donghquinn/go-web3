package web3

import (
	"math/big"
	"strings"
	"testing"
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

	// Compare by hex representation
	if PrivateKeyToHex(recovered) != hexKey {
		t.Error("round-trip private key mismatch")
	}
}

func TestPrivateKeyFromHexWithout0x(t *testing.T) {
	key, _ := GeneratePrivateKey()
	hexKey := PrivateKeyToHex(key)
	// strip 0x prefix
	stripped := hexKey[2:]
	recovered, err := PrivateKeyFromHex(stripped)
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
	// Known private key → known address
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

func TestSignTransactionMissingTo(t *testing.T) {
	key, _ := GeneratePrivateKey()
	tp := NewTransactionParams().
		SetGas(21000).
		SetGasPrice(big.NewInt(1))
	_, err := SignTransaction(tp, key)
	if err == nil {
		t.Error("expected error when To is missing")
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
	tp.Nonce = 0
	tp.ChainID = ChainMainnet.BigInt()

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

func TestRecoverSignerInvalidHex(t *testing.T) {
	_, err := RecoverSigner("0xnotvalidhex")
	if err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestCreateContractCall(t *testing.T) {
	key, _ := GeneratePrivateKey()
	params := NewTransactionParams().
		SetGas(100000).
		SetGasPrice(big.NewInt(20000000000)).
		SetNonce(0).
		SetChainID(ChainMainnet)

	methodData := []byte{0xde, 0xad, 0xbe, 0xef}
	contract := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"

	signed, err := CreateContractCall(contract, methodData, key, params)
	if err != nil {
		t.Fatalf("CreateContractCall() error = %v", err)
	}
	if !strings.HasPrefix(signed.Hash, "0x") {
		t.Errorf("Hash = %q, want 0x prefix", signed.Hash)
	}
}

func TestRandomNonce(t *testing.T) {
	n1 := RandomNonce()
	n2 := RandomNonce()
	// Both should be valid uint64; they're almost certainly different
	_ = n1
	_ = n2
	// Just verify the function runs without panic
}
