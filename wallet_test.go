package web3

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// walletMockServer routes RPC methods to responses for wallet tests.
func walletMockServer(t *testing.T, handlers map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req RPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		result, ok := handlers[req.Method]
		if !ok {
			result = "0x0"
		}

		resultBytes, _ := json.Marshal(result)
		resp := RPCResponse{
			ID:     req.ID,
			Result: json.RawMessage(resultBytes),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

const testPrivKeyHex = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
const testAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

func TestNewWallet(t *testing.T) {
	c := NewClient("http://localhost:8545")
	w, err := NewWallet(testPrivKeyHex, c)
	if err != nil {
		t.Fatalf("NewWallet() error = %v", err)
	}
	if !strings.EqualFold(w.GetAddress(), testAddress) {
		t.Errorf("GetAddress() = %q, want %q", w.GetAddress(), testAddress)
	}
}

func TestNewWalletInvalidKey(t *testing.T) {
	c := NewClient("http://localhost:8545")
	_, err := NewWallet("0xinvalidkey", c)
	if err == nil {
		t.Error("expected error for invalid private key")
	}
}

func TestCreateWallet(t *testing.T) {
	c := NewClient("http://localhost:8545")
	w, err := CreateWallet(c)
	if err != nil {
		t.Fatalf("CreateWallet() error = %v", err)
	}
	if !IsAddress(w.GetAddress()) {
		t.Errorf("CreateWallet() address %q is not valid", w.GetAddress())
	}
}

func TestWalletGetPrivateKey(t *testing.T) {
	c := NewClient("http://localhost:8545")
	w, _ := NewWallet(testPrivKeyHex, c)
	key := w.GetPrivateKey()
	if !strings.HasPrefix(key, "0x") {
		t.Errorf("GetPrivateKey() = %q, want 0x prefix", key)
	}
}

func TestWalletGetBalance(t *testing.T) {
	srv := walletMockServer(t, map[string]interface{}{
		"eth_getBalance": "0xde0b6b3a7640000", // 1 ETH
	})
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
	srv := walletMockServer(t, map[string]interface{}{
		"eth_getTransactionCount": "0x5",
	})
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
	srv := walletMockServer(t, map[string]interface{}{
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
	srv := walletMockServer(t, map[string]interface{}{
		"eth_gasPrice":            "0x4a817c800",
		"eth_getTransactionCount": "0x1",
		"eth_sendRawTransaction":  "0xtxhash2",
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	result, err := w.SendTransaction(context.Background(), &TransferOptions{
		To:       "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		Value:    big.NewInt(0),
		GasLimit: 21000, // pre-set, skips estimation
	})
	if err != nil {
		t.Fatalf("SendTransaction() with preset gas error = %v", err)
	}
	if result.TransactionHash == "" {
		t.Error("TransactionHash is empty")
	}
}

func TestWalletSendEther(t *testing.T) {
	srv := walletMockServer(t, map[string]interface{}{
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
	c := NewClient("http://localhost:8545")
	w, _ := NewWallet(testPrivKeyHex, c)
	_, err := w.SendEther(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", "not-a-number")
	if err == nil {
		t.Error("expected error for invalid ether amount")
	}
}

func TestWalletSendWei(t *testing.T) {
	srv := walletMockServer(t, map[string]interface{}{
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
	srv := walletMockServer(t, map[string]interface{}{
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
		big.NewInt(30000000000), // maxFeePerGas
		big.NewInt(2000000000),  // maxPriorityFeePerGas
	)
	if err != nil {
		t.Fatalf("SendEIP1559Transaction() error = %v", err)
	}
	if result.TransactionHash != "0x1559hash" {
		t.Errorf("TransactionHash = %q", result.TransactionHash)
	}
}

func TestWalletCallContract(t *testing.T) {
	srv := walletMockServer(t, map[string]interface{}{
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
	srv := walletMockServer(t, map[string]interface{}{
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

func TestWalletWaitForTransaction(t *testing.T) {
	receipt := TransactionReceipt{
		TransactionHash: "0xhash",
		Status:          "0x1",
	}
	srv := walletMockServer(t, map[string]interface{}{
		"eth_getTransactionReceipt": receipt,
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
	// Server that always returns null (no receipt yet)
	srv := walletMockServer(t, map[string]interface{}{
		"eth_getTransactionReceipt": nil,
	})
	defer srv.Close()

	w, _ := NewWallet(testPrivKeyHex, NewClient(srv.URL))
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // immediately cancel

	_, err := w.WaitForTransaction(ctx, "0xhash")
	if err == nil {
		t.Error("expected context cancellation error")
	}
}
