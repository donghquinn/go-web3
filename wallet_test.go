package web3

import (
	"context"
	"math/big"
	"strings"
	"testing"
)

const testPrivKeyHex = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
const testAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

func TestNewWallet(t *testing.T) {
	w, err := NewWallet(testPrivKeyHex, NewClient("http://localhost:8545"))
	if err != nil {
		t.Fatalf("NewWallet() error = %v", err)
	}
	if !strings.EqualFold(w.GetAddress(), testAddress) {
		t.Errorf("GetAddress() = %q, want %q", w.GetAddress(), testAddress)
	}
}

func TestNewWalletInvalidKey(t *testing.T) {
	_, err := NewWallet("0xinvalidkey", NewClient("http://localhost:8545"))
	if err == nil {
		t.Error("expected error for invalid private key")
	}
}

func TestCreateWallet(t *testing.T) {
	w, err := CreateWallet(NewClient("http://localhost:8545"))
	if err != nil {
		t.Fatalf("CreateWallet() error = %v", err)
	}
	if !IsAddress(w.GetAddress()) {
		t.Errorf("CreateWallet() address %q is not valid", w.GetAddress())
	}
}

func TestWalletGetPrivateKey(t *testing.T) {
	w, _ := NewWallet(testPrivKeyHex, NewClient("http://localhost:8545"))
	if key := w.GetPrivateKey(); !strings.HasPrefix(key, "0x") {
		t.Errorf("GetPrivateKey() = %q, want 0x prefix", key)
	}
}

func TestWalletGetBalance(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{"eth_getBalance": "0xde0b6b3a7640000"})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	balance, err := w.GetBalance(context.Background())
	if err != nil {
		t.Fatalf("GetBalance() error = %v", err)
	}
	expected, _ := ToWei("1", Ether)
	if balance.Cmp(expected) != 0 {
		t.Errorf("GetBalance() = %v, want %v", balance, expected)
	}
}

func TestWalletGetNonce(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{"eth_getTransactionCount": "0x5"})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	nonce, err := w.GetNonce(context.Background())
	if err != nil {
		t.Fatalf("GetNonce() error = %v", err)
	}
	if nonce != 5 {
		t.Errorf("GetNonce() = %d, want 5", nonce)
	}
}

func TestWalletSendTransaction(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		"eth_estimateGas":         "0x5208",
		"eth_gasPrice":            "0x4a817c800",
		"eth_getTransactionCount": "0x0",
		"eth_sendRawTransaction":  "0xtxhash",
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	result, err := w.SendTransaction(context.Background(), &TransferOptions{
		To:    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		Value: big.NewInt(0),
	})
	if err != nil {
		t.Fatalf("SendTransaction() error = %v", err)
	}
	if result.TransactionHash != "0xtxhash" {
		t.Errorf("TransactionHash = %q, want %q", result.TransactionHash, "0xtxhash")
	}
	if !strings.EqualFold(result.From, testAddress) {
		t.Errorf("From = %q, want %q", result.From, testAddress)
	}
}

func TestWalletSendTransactionWithPresetGas(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		"eth_gasPrice":            "0x4a817c800",
		"eth_getTransactionCount": "0x1",
		"eth_sendRawTransaction":  "0xtxhash2",
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	result, err := w.SendTransaction(context.Background(), &TransferOptions{
		To:       "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		Value:    big.NewInt(0),
		GasLimit: 21000,
	})
	if err != nil {
		t.Fatalf("SendTransaction() with preset gas error = %v", err)
	}
	if result.TransactionHash == "" {
		t.Error("TransactionHash is empty")
	}
}

func TestWalletSendTransactionGasEstimateError(t *testing.T) {
	srv := newMockRPCErrorServer(-32000, "estimate error")
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	_, err := w.SendTransaction(context.Background(), &TransferOptions{
		To:    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		Value: big.NewInt(0),
	})
	if err == nil {
		t.Fatal("expected error from gas estimation failure")
	}
}

func TestWalletSendTransactionGasPriceError(t *testing.T) {
	// GasLimit pre-set (skips estimation); all RPC calls fail → gas price fetch fails.
	srv := newMockRPCErrorServer(-32000, "rpc error")
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	_, err := w.SendTransaction(context.Background(), &TransferOptions{
		To:       "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		Value:    big.NewInt(0),
		GasLimit: 21000,
	})
	if err == nil {
		t.Fatal("expected error from gas price fetch failure")
	}
}

func TestWalletSendTransactionNonceError(t *testing.T) {
	// Gas resolved but nonce fetch fails.
	srv := newMockRPCErrorServer(-32000, "nonce error")
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	_, err := w.SendTransaction(context.Background(), &TransferOptions{
		To:       "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		Value:    big.NewInt(0),
		GasLimit: 21000,
		GasPrice: big.NewInt(1),
	})
	if err == nil {
		t.Fatal("expected error from nonce fetch failure")
	}
}

func TestWalletSendEther(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		"eth_estimateGas":         "0x5208",
		"eth_gasPrice":            "0x4a817c800",
		"eth_getTransactionCount": "0x0",
		"eth_sendRawTransaction":  "0xethhash",
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	result, err := w.SendEther(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", "0.001")
	if err != nil {
		t.Fatalf("SendEther() error = %v", err)
	}
	if result.TransactionHash != "0xethhash" {
		t.Errorf("TransactionHash = %q", result.TransactionHash)
	}
}

func TestWalletSendEtherInvalidAmount(t *testing.T) {
	w, _ := NewWallet(testPrivKeyHex, NewClient("http://localhost:8545"))
	_, err := w.SendEther(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", "not-a-number")
	if err == nil {
		t.Error("expected error for invalid ether amount")
	}
}

func TestWalletSendWei(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		"eth_estimateGas":         "0x5208",
		"eth_gasPrice":            "0x4a817c800",
		"eth_getTransactionCount": "0x0",
		"eth_sendRawTransaction":  "0xweihash",
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	result, err := w.SendWei(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", big.NewInt(1000))
	if err != nil {
		t.Fatalf("SendWei() error = %v", err)
	}
	if result.TransactionHash == "" {
		t.Error("SendWei() TransactionHash is empty")
	}
}

func TestWalletSendEIP1559Transaction(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		"eth_estimateGas":         "0x5208",
		"eth_getTransactionCount": "0x0",
		"eth_sendRawTransaction":  "0x1559hash",
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	result, err := w.SendEIP1559Transaction(
		context.Background(),
		&TransferOptions{
			To:    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			Value: big.NewInt(0),
		},
		big.NewInt(30000000000),
		big.NewInt(2000000000),
	)
	if err != nil {
		t.Fatalf("SendEIP1559Transaction() error = %v", err)
	}
	if result.TransactionHash != "0x1559hash" {
		t.Errorf("TransactionHash = %q", result.TransactionHash)
	}
}

func TestWalletSendTransactionSendError(t *testing.T) {
	// All prerequisite RPCs succeed but sendRawTransaction returns a non-string → unmarshal error.
	srv := newDispatchServer(map[string]interface{}{
		"eth_estimateGas":         "0x5208",
		"eth_gasPrice":            "0x4a817c800",
		"eth_getTransactionCount": "0x0",
		"eth_sendRawTransaction":  42, // int, not string → parseJSONString fails
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	_, err := w.SendTransaction(context.Background(), &TransferOptions{
		To:    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		Value: big.NewInt(0),
	})
	if err == nil {
		t.Fatal("expected send transaction error")
	}
}

func TestWalletSendEIP1559TransactionGasEstimateError(t *testing.T) {
	srv := newMockRPCErrorServer(-32000, "estimate error")
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	_, err := w.SendEIP1559Transaction(
		context.Background(),
		&TransferOptions{To: "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", Value: big.NewInt(0)},
		big.NewInt(1), big.NewInt(1),
	)
	if err == nil {
		t.Fatal("expected gas estimation error")
	}
}

func TestWalletSendEIP1559TransactionNonceError(t *testing.T) {
	srv := newMockRPCErrorServer(-32000, "nonce error")
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	_, err := w.SendEIP1559Transaction(
		context.Background(),
		&TransferOptions{To: "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", Value: big.NewInt(0), GasLimit: 21000},
		big.NewInt(1), big.NewInt(1),
	)
	if err == nil {
		t.Fatal("expected nonce fetch error")
	}
}

func TestWalletSendEIP1559TransactionSignError(t *testing.T) {
	// Empty To triggers SignEIP1559Transaction to return "recipient is required" error.
	srv := newDispatchServer(map[string]interface{}{
		"eth_getTransactionCount": "0x0",
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	_, err := w.SendEIP1559Transaction(
		context.Background(),
		&TransferOptions{To: "", Value: big.NewInt(0), GasLimit: 21000},
		big.NewInt(30000000000),
		big.NewInt(2000000000),
	)
	if err == nil {
		t.Fatal("expected SignEIP1559Transaction error for empty recipient")
	}
}

func TestWalletSendEIP1559TransactionSendError(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		"eth_estimateGas":         "0x5208",
		"eth_getTransactionCount": "0x0",
		"eth_sendRawTransaction":  42, // int → parseJSONString fails
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	_, err := w.SendEIP1559Transaction(
		context.Background(),
		&TransferOptions{To: "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", Value: big.NewInt(0)},
		big.NewInt(30000000000),
		big.NewInt(2000000000),
	)
	if err == nil {
		t.Fatal("expected send EIP-1559 transaction error")
	}
}

func TestWalletDeployContractGasPriceError(t *testing.T) {
	// GasLimit preset so estimation is skipped; gasPrice fetch fails.
	srv := newMockRPCErrorServer(-32000, "rpc error")
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	_, err := w.DeployContract(context.Background(), []byte{0x60}, nil, 500000, nil)
	if err == nil {
		t.Fatal("expected gas price fetch error")
	}
}

func TestWalletCallContract(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		"eth_call": "0x000000000000000000000000000000000000000000000000000000000000002a",
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	result, err := w.CallContract(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", []byte{0x70, 0xa0, 0x82, 0x31})
	if err != nil {
		t.Fatalf("CallContract() error = %v", err)
	}
	if result == "" {
		t.Error("CallContract() returned empty result")
	}
}

func TestWalletSendContractTransaction(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		"eth_estimateGas":         "0x186a0",
		"eth_gasPrice":            "0x4a817c800",
		"eth_getTransactionCount": "0x0",
		"eth_sendRawTransaction":  "0xcontracthash",
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	result, err := w.SendContractTransaction(
		context.Background(),
		"0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		[]byte{0xde, 0xad},
		big.NewInt(0),
	)
	if err != nil {
		t.Fatalf("SendContractTransaction() error = %v", err)
	}
	if result.TransactionHash == "" {
		t.Error("TransactionHash is empty")
	}
}

func TestWalletDeployContract(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		"eth_estimateGas":         "0x186a0",
		"eth_gasPrice":            "0x4a817c800",
		"eth_getTransactionCount": "0x0",
		"eth_sendRawTransaction":  "0xdeployhash",
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	result, err := w.DeployContract(
		context.Background(),
		[]byte{0x60, 0x80, 0x60, 0x40},
		[]byte{0x00, 0x01},
		0,   // gasLimit=0 → triggers estimation
		nil, // gasPrice=nil → triggers fetch
	)
	if err != nil {
		t.Fatalf("DeployContract() error = %v", err)
	}
	if result.TransactionHash == "" {
		t.Error("DeployContract() TransactionHash is empty")
	}
}

func TestWalletDeployContractWithPresetParams(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		"eth_getTransactionCount": "0x0",
		"eth_sendRawTransaction":  "0xdeployhash2",
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	result, err := w.DeployContract(
		context.Background(),
		[]byte{0x60, 0x80},
		nil,
		500000,                  // preset gas limit
		big.NewInt(20000000000), // preset gas price
	)
	if err != nil {
		t.Fatalf("DeployContract() with preset params error = %v", err)
	}
	if result.TransactionHash == "" {
		t.Error("DeployContract() TransactionHash is empty")
	}
}

func TestWalletDeployContractGasEstimateError(t *testing.T) {
	srv := newMockRPCErrorServer(-32000, "estimate error")
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	_, err := w.DeployContract(context.Background(), []byte{0x60}, nil, 0, nil)
	if err == nil {
		t.Fatal("expected gas estimation error")
	}
}

func TestWalletWaitForTransaction(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		"eth_getTransactionReceipt": TransactionReceipt{TransactionHash: "0xhash", Status: "0x1"},
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	got, err := w.WaitForTransaction(context.Background(), "0xhash")
	if err != nil {
		t.Fatalf("WaitForTransaction() error = %v", err)
	}
	if got.Status != "0x1" {
		t.Errorf("Status = %q, want %q", got.Status, "0x1")
	}
}

func TestWalletWaitForTransactionContextCancel(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{"eth_getTransactionReceipt": nil})
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	_, err := w.WaitForTransaction(ctx, "0xhash")
	if err == nil {
		t.Error("expected context cancellation error")
	}
}
