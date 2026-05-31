package web3

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// dispatchMockServer routes different RPC methods to different result values.
func dispatchMockServer(handlers map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req RPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		result, ok := handlers[req.Method]
		if !ok {
			result = nil
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

func TestGetBalance(t *testing.T) {
	srv := newMockRPCServer("0xde0b6b3a7640000") // 1 ETH in wei
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	balance, err := eth.GetBalance(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", BlockLatest)
	if err != nil {
		t.Fatalf("GetBalance() error = %v", err)
	}
	expected, _ := ToWei("1", Ether)
	if balance.Cmp(expected) != 0 {
		t.Errorf("GetBalance() = %v, want %v", balance, expected)
	}
}

func TestGetBalanceDefaultBlock(t *testing.T) {
	srv := newMockRPCServer("0x0")
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	_, err := eth.GetBalance(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", "")
	if err != nil {
		t.Fatalf("GetBalance() with empty block = %v", err)
	}
}

func TestGetBlockNumber(t *testing.T) {
	srv := newMockRPCServer("0x100") // block 256
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	num, err := eth.GetBlockNumber(context.Background())
	if err != nil {
		t.Fatalf("GetBlockNumber() error = %v", err)
	}
	if num != 256 {
		t.Errorf("GetBlockNumber() = %d, want 256", num)
	}
}

func TestGetGasPrice(t *testing.T) {
	srv := newMockRPCServer("0x4a817c800") // 20 gwei
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	price, err := eth.GetGasPrice(context.Background())
	if err != nil {
		t.Fatalf("GetGasPrice() error = %v", err)
	}
	expected, _ := ToWei("20", Gwei)
	if price.Cmp(expected) != 0 {
		t.Errorf("GetGasPrice() = %v, want %v", price, expected)
	}
}

func TestGetTransactionCount(t *testing.T) {
	srv := newMockRPCServer("0x5") // nonce 5
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	count, err := eth.GetTransactionCount(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", BlockLatest)
	if err != nil {
		t.Fatalf("GetTransactionCount() error = %v", err)
	}
	if count != 5 {
		t.Errorf("GetTransactionCount() = %d, want 5", count)
	}
}

func TestGetChainID(t *testing.T) {
	srv := newMockRPCServer("0x1") // mainnet
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	chainID, err := eth.GetChainID(context.Background())
	if err != nil {
		t.Fatalf("GetChainID() error = %v", err)
	}
	if chainID.Int64() != 1 {
		t.Errorf("GetChainID() = %v, want 1", chainID)
	}
}

func TestGetNetVersion(t *testing.T) {
	srv := newMockRPCServer("1")
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	ver, err := eth.GetNetVersion(context.Background())
	if err != nil {
		t.Fatalf("GetNetVersion() error = %v", err)
	}
	if ver != "1" {
		t.Errorf("GetNetVersion() = %q, want %q", ver, "1")
	}
}

func TestGetClientVersion(t *testing.T) {
	srv := newMockRPCServer("Geth/v1.11.0")
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	ver, err := eth.GetClientVersion(context.Background())
	if err != nil {
		t.Fatalf("GetClientVersion() error = %v", err)
	}
	if ver != "Geth/v1.11.0" {
		t.Errorf("GetClientVersion() = %q", ver)
	}
}

func TestGetMaxPriorityFeePerGas(t *testing.T) {
	srv := newMockRPCServer("0x77359400") // 2 gwei
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	fee, err := eth.GetMaxPriorityFeePerGas(context.Background())
	if err != nil {
		t.Fatalf("GetMaxPriorityFeePerGas() error = %v", err)
	}
	expected, _ := ToWei("2", Gwei)
	if fee.Cmp(expected) != 0 {
		t.Errorf("GetMaxPriorityFeePerGas() = %v, want %v", fee, expected)
	}
}

func TestGetBlockByNumber(t *testing.T) {
	block := Block{
		Number:       "0x1",
		Hash:         "0xabc",
		Transactions: []interface{}{},
	}
	srv := newMockRPCServer(block)
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	b, err := eth.GetBlockByNumber(context.Background(), BlockLatest, false)
	if err != nil {
		t.Fatalf("GetBlockByNumber() error = %v", err)
	}
	if b.Hash != "0xabc" {
		t.Errorf("GetBlockByNumber() Hash = %q, want %q", b.Hash, "0xabc")
	}
}

func TestGetBlockByNumberDefaultParam(t *testing.T) {
	block := Block{Number: "0x1", Transactions: []interface{}{}}
	srv := newMockRPCServer(block)
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	_, err := eth.GetBlockByNumber(context.Background(), "", false)
	if err != nil {
		t.Fatalf("GetBlockByNumber() with empty block param error = %v", err)
	}
}

func TestGetBlockByHash(t *testing.T) {
	block := Block{
		Hash:         "0xdeadbeef",
		Number:       "0x5",
		Transactions: []interface{}{},
	}
	srv := newMockRPCServer(block)
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	b, err := eth.GetBlockByHash(context.Background(), "0xdeadbeef", false)
	if err != nil {
		t.Fatalf("GetBlockByHash() error = %v", err)
	}
	if b.Hash != "0xdeadbeef" {
		t.Errorf("GetBlockByHash() Hash = %q", b.Hash)
	}
}

func TestGetTransactionByHash(t *testing.T) {
	tx := Transaction{
		Hash:  "0xabc123",
		From:  "0x1234",
		To:    "0x5678",
		Value: "0x0",
	}
	srv := newMockRPCServer(tx)
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	got, err := eth.GetTransactionByHash(context.Background(), "0xabc123")
	if err != nil {
		t.Fatalf("GetTransactionByHash() error = %v", err)
	}
	if got.Hash != "0xabc123" {
		t.Errorf("GetTransactionByHash() Hash = %q", got.Hash)
	}
}

func TestGetTransactionReceipt(t *testing.T) {
	receipt := TransactionReceipt{
		TransactionHash: "0xabc123",
		Status:          "0x1",
		BlockNumber:     "0xa",
		GasUsed:         "0x5208",
	}
	srv := newMockRPCServer(receipt)
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	got, err := eth.GetTransactionReceipt(context.Background(), "0xabc123")
	if err != nil {
		t.Fatalf("GetTransactionReceipt() error = %v", err)
	}
	if got.Status != "0x1" {
		t.Errorf("GetTransactionReceipt() Status = %q, want %q", got.Status, "0x1")
	}
}

func TestSendRawTransaction(t *testing.T) {
	srv := newMockRPCServer("0xhash123")
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	hash, err := eth.SendRawTransaction(context.Background(), "0xrawdata")
	if err != nil {
		t.Fatalf("SendRawTransaction() error = %v", err)
	}
	if hash != "0xhash123" {
		t.Errorf("SendRawTransaction() = %q, want %q", hash, "0xhash123")
	}
}

func TestEstimateGas(t *testing.T) {
	srv := newMockRPCServer("0x5208") // 21000
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	gas, err := eth.EstimateGas(context.Background(), map[string]interface{}{
		"to":    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		"value": "0x0",
	})
	if err != nil {
		t.Fatalf("EstimateGas() error = %v", err)
	}
	if gas != 21000 {
		t.Errorf("EstimateGas() = %d, want 21000", gas)
	}
}

func TestCall(t *testing.T) {
	srv := newMockRPCServer("0x000000000000000000000000000000000000000000000000000000000000002a")
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	result, err := eth.Call(context.Background(), map[string]interface{}{
		"to":   "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		"data": "0x70a08231",
	}, BlockLatest)
	if err != nil {
		t.Fatalf("Call() error = %v", err)
	}
	if result == "" {
		t.Error("Call() returned empty result")
	}
}

func TestCallDefaultBlock(t *testing.T) {
	srv := newMockRPCServer("0x0")
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	_, err := eth.Call(context.Background(), map[string]interface{}{"to": "0x1"}, "")
	if err != nil {
		t.Fatalf("Call() with empty block error = %v", err)
	}
}

func TestGetStorageAt(t *testing.T) {
	srv := newMockRPCServer("0x000000000000000000000000000000000000000000000000000000000000002a")
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	val, err := eth.GetStorageAt(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", "0x0", BlockLatest)
	if err != nil {
		t.Fatalf("GetStorageAt() error = %v", err)
	}
	if val == "" {
		t.Error("GetStorageAt() returned empty")
	}
}

func TestGetStorageAtDefaultBlock(t *testing.T) {
	srv := newMockRPCServer("0x0")
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	_, err := eth.GetStorageAt(context.Background(), "0x1", "0x0", "")
	if err != nil {
		t.Fatalf("GetStorageAt() with empty block = %v", err)
	}
}

func TestGetCode(t *testing.T) {
	srv := newMockRPCServer("0x6080604052")
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	code, err := eth.GetCode(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", BlockLatest)
	if err != nil {
		t.Fatalf("GetCode() error = %v", err)
	}
	if code != "0x6080604052" {
		t.Errorf("GetCode() = %q", code)
	}
}

func TestGetLogs(t *testing.T) {
	logs := []*Log{
		{
			Address:         "0xabc",
			TransactionHash: "0xdef",
			BlockNumber:     "0x1",
		},
	}
	srv := newMockRPCServer(logs)
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	got, err := eth.GetLogs(context.Background(), LogFilter{
		FromBlock: BlockLatest,
		ToBlock:   BlockLatest,
	})
	if err != nil {
		t.Fatalf("GetLogs() error = %v", err)
	}
	if len(got) != 1 {
		t.Errorf("GetLogs() returned %d logs, want 1", len(got))
	}
}

func TestGetFeeHistory(t *testing.T) {
	feeHistory := FeeHistoryResult{
		OldestBlock:   "0x1",
		BaseFeePerGas: []string{"0x4a817c800"},
		GasUsedRatio:  []float64{0.5},
	}
	srv := newMockRPCServer(feeHistory)
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	got, err := eth.GetFeeHistory(context.Background(), 1, BlockLatest, []float64{25, 50, 75})
	if err != nil {
		t.Fatalf("GetFeeHistory() error = %v", err)
	}
	if len(got.BaseFeePerGas) == 0 {
		t.Error("GetFeeHistory() returned empty BaseFeePerGas")
	}
}

func TestGetFeeHistoryDefaultBlock(t *testing.T) {
	feeHistory := FeeHistoryResult{OldestBlock: "0x1"}
	srv := newMockRPCServer(feeHistory)
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	_, err := eth.GetFeeHistory(context.Background(), 1, "", nil)
	if err != nil {
		t.Fatalf("GetFeeHistory() with empty block = %v", err)
	}
}

func TestGetPendingTransactions(t *testing.T) {
	block := map[string]interface{}{
		"number": "pending",
		"hash":   "0x0",
		"transactions": []interface{}{
			map[string]interface{}{
				"hash":     "0xpending1",
				"from":     "0xsender",
				"to":       "0xrecipient",
				"value":    "0x0",
				"gas":      "0x5208",
				"gasPrice": "0x4a817c800",
				"nonce":    "0x1",
				"input":    "0x",
			},
		},
	}
	srv := newMockRPCServer(block)
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	txs, err := eth.GetPendingTransactions(context.Background())
	if err != nil {
		t.Fatalf("GetPendingTransactions() error = %v", err)
	}
	if len(txs) != 1 {
		t.Errorf("GetPendingTransactions() returned %d txs, want 1", len(txs))
	}
	if txs[0].Hash != "0xpending1" {
		t.Errorf("tx Hash = %q, want %q", txs[0].Hash, "0xpending1")
	}
}

func TestGetPendingTransactionCount(t *testing.T) {
	block := map[string]interface{}{
		"number":       "pending",
		"hash":         "0x0",
		"transactions": []interface{}{},
	}
	srv := newMockRPCServer(block)
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	count, err := eth.GetPendingTransactionCount(context.Background())
	if err != nil {
		t.Fatalf("GetPendingTransactionCount() error = %v", err)
	}
	if count != 0 {
		t.Errorf("GetPendingTransactionCount() = %d, want 0", count)
	}
}

func TestGetAccountPendingTransactions(t *testing.T) {
	const addr = "0xsender"
	block := map[string]interface{}{
		"number": "pending",
		"hash":   "0x0",
		"transactions": []interface{}{
			map[string]interface{}{
				"hash":  "0xtx1",
				"from":  addr,
				"to":    "0xother",
				"value": "0x0",
			},
			map[string]interface{}{
				"hash":  "0xtx2",
				"from":  "0xother",
				"to":    "0xunrelated",
				"value": "0x0",
			},
		},
	}
	srv := newMockRPCServer(block)
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	txs, err := eth.GetAccountPendingTransactions(context.Background(), addr)
	if err != nil {
		t.Fatalf("GetAccountPendingTransactions() error = %v", err)
	}
	if len(txs) != 1 {
		t.Errorf("GetAccountPendingTransactions() returned %d txs, want 1", len(txs))
	}
}

func TestIsPendingTransaction(t *testing.T) {
	block := map[string]interface{}{
		"number": "pending",
		"hash":   "0x0",
		"transactions": []interface{}{
			map[string]interface{}{
				"hash":  "0xtarget",
				"from":  "0xa",
				"to":    "0xb",
				"value": "0x0",
			},
		},
	}
	srv := newMockRPCServer(block)
	defer srv.Close()

	eth := NewClient(srv.URL).Eth()
	isPending, err := eth.IsPendingTransaction(context.Background(), "0xtarget")
	if err != nil {
		t.Fatalf("IsPendingTransaction() error = %v", err)
	}
	if !isPending {
		t.Error("IsPendingTransaction() = false, want true")
	}

	notPending, err := eth.IsPendingTransaction(context.Background(), "0xnothere")
	if err != nil {
		t.Fatalf("IsPendingTransaction() error = %v", err)
	}
	if notPending {
		t.Error("IsPendingTransaction() = true for absent tx, want false")
	}
}
