package web3

import (
	"context"
	"testing"
)

func TestGetBalance(t *testing.T) {
	srv := newMockRPCServer("0xde0b6b3a7640000") // 1 ETH in wei
	defer srv.Close()

	balance, err := NewClient(srv.URL).Eth().GetBalance(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", BlockLatest)
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

	if _, err := NewClient(srv.URL).Eth().GetBalance(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", ""); err != nil {
		t.Fatalf("GetBalance() with empty block = %v", err)
	}
}

func TestGetBlockNumber(t *testing.T) {
	srv := newMockRPCServer("0x100") // block 256
	defer srv.Close()

	num, err := NewClient(srv.URL).Eth().GetBlockNumber(context.Background())
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

	price, err := NewClient(srv.URL).Eth().GetGasPrice(context.Background())
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

	count, err := NewClient(srv.URL).Eth().GetTransactionCount(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", BlockLatest)
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

	chainID, err := NewClient(srv.URL).Eth().GetChainID(context.Background())
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

	ver, err := NewClient(srv.URL).Eth().GetNetVersion(context.Background())
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

	ver, err := NewClient(srv.URL).Eth().GetClientVersion(context.Background())
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

	fee, err := NewClient(srv.URL).Eth().GetMaxPriorityFeePerGas(context.Background())
	if err != nil {
		t.Fatalf("GetMaxPriorityFeePerGas() error = %v", err)
	}
	expected, _ := ToWei("2", Gwei)
	if fee.Cmp(expected) != 0 {
		t.Errorf("GetMaxPriorityFeePerGas() = %v, want %v", fee, expected)
	}
}

func TestGetBlockByNumber(t *testing.T) {
	srv := newMockRPCServer(Block{Number: "0x1", Hash: "0xabc", Transactions: []interface{}{}})
	defer srv.Close()

	b, err := NewClient(srv.URL).Eth().GetBlockByNumber(context.Background(), BlockLatest, false)
	if err != nil {
		t.Fatalf("GetBlockByNumber() error = %v", err)
	}
	if b.Hash != "0xabc" {
		t.Errorf("GetBlockByNumber() Hash = %q, want %q", b.Hash, "0xabc")
	}
}

func TestGetBlockByNumberDefaultParam(t *testing.T) {
	srv := newMockRPCServer(Block{Number: "0x1", Transactions: []interface{}{}})
	defer srv.Close()

	if _, err := NewClient(srv.URL).Eth().GetBlockByNumber(context.Background(), "", false); err != nil {
		t.Fatalf("GetBlockByNumber() with empty block param error = %v", err)
	}
}

func TestGetBlockByHash(t *testing.T) {
	srv := newMockRPCServer(Block{Hash: "0xdeadbeef", Number: "0x5", Transactions: []interface{}{}})
	defer srv.Close()

	b, err := NewClient(srv.URL).Eth().GetBlockByHash(context.Background(), "0xdeadbeef", false)
	if err != nil {
		t.Fatalf("GetBlockByHash() error = %v", err)
	}
	if b.Hash != "0xdeadbeef" {
		t.Errorf("GetBlockByHash() Hash = %q", b.Hash)
	}
}

func TestGetTransactionByHash(t *testing.T) {
	srv := newMockRPCServer(Transaction{Hash: "0xabc123", From: "0x1234", To: "0x5678", Value: "0x0"})
	defer srv.Close()

	tx, err := NewClient(srv.URL).Eth().GetTransactionByHash(context.Background(), "0xabc123")
	if err != nil {
		t.Fatalf("GetTransactionByHash() error = %v", err)
	}
	if tx.Hash != "0xabc123" {
		t.Errorf("GetTransactionByHash() Hash = %q", tx.Hash)
	}
}

func TestGetTransactionReceipt(t *testing.T) {
	srv := newMockRPCServer(TransactionReceipt{TransactionHash: "0xabc123", Status: "0x1", BlockNumber: "0xa", GasUsed: "0x5208"})
	defer srv.Close()

	receipt, err := NewClient(srv.URL).Eth().GetTransactionReceipt(context.Background(), "0xabc123")
	if err != nil {
		t.Fatalf("GetTransactionReceipt() error = %v", err)
	}
	if receipt.Status != "0x1" {
		t.Errorf("GetTransactionReceipt() Status = %q, want %q", receipt.Status, "0x1")
	}
}

func TestSendRawTransaction(t *testing.T) {
	srv := newMockRPCServer("0xhash123")
	defer srv.Close()

	hash, err := NewClient(srv.URL).Eth().SendRawTransaction(context.Background(), "0xrawdata")
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

	gas, err := NewClient(srv.URL).Eth().EstimateGas(context.Background(), map[string]interface{}{
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

	result, err := NewClient(srv.URL).Eth().Call(context.Background(), map[string]interface{}{
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

	if _, err := NewClient(srv.URL).Eth().Call(context.Background(), map[string]interface{}{"to": "0x1"}, ""); err != nil {
		t.Fatalf("Call() with empty block error = %v", err)
	}
}

func TestGetStorageAt(t *testing.T) {
	srv := newMockRPCServer("0x000000000000000000000000000000000000000000000000000000000000002a")
	defer srv.Close()

	val, err := NewClient(srv.URL).Eth().GetStorageAt(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", "0x0", BlockLatest)
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

	if _, err := NewClient(srv.URL).Eth().GetStorageAt(context.Background(), "0x1", "0x0", ""); err != nil {
		t.Fatalf("GetStorageAt() with empty block = %v", err)
	}
}

func TestGetCode(t *testing.T) {
	srv := newMockRPCServer("0x6080604052")
	defer srv.Close()

	code, err := NewClient(srv.URL).Eth().GetCode(context.Background(), "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", BlockLatest)
	if err != nil {
		t.Fatalf("GetCode() error = %v", err)
	}
	if code != "0x6080604052" {
		t.Errorf("GetCode() = %q", code)
	}
}

func TestGetLogs(t *testing.T) {
	logs := []*Log{{Address: "0xabc", TransactionHash: "0xdef", BlockNumber: "0x1"}}
	srv := newMockRPCServer(logs)
	defer srv.Close()

	got, err := NewClient(srv.URL).Eth().GetLogs(context.Background(), LogFilter{FromBlock: BlockLatest, ToBlock: BlockLatest})
	if err != nil {
		t.Fatalf("GetLogs() error = %v", err)
	}
	if len(got) != 1 {
		t.Errorf("GetLogs() returned %d logs, want 1", len(got))
	}
}

func TestGetFeeHistory(t *testing.T) {
	srv := newMockRPCServer(FeeHistoryResult{
		OldestBlock:   "0x1",
		BaseFeePerGas: []string{"0x4a817c800"},
		GasUsedRatio:  []float64{0.5},
	})
	defer srv.Close()

	got, err := NewClient(srv.URL).Eth().GetFeeHistory(context.Background(), 1, BlockLatest, []float64{25, 50, 75})
	if err != nil {
		t.Fatalf("GetFeeHistory() error = %v", err)
	}
	if len(got.BaseFeePerGas) == 0 {
		t.Error("GetFeeHistory() returned empty BaseFeePerGas")
	}
}

func TestGetFeeHistoryDefaultBlock(t *testing.T) {
	srv := newMockRPCServer(FeeHistoryResult{OldestBlock: "0x1"})
	defer srv.Close()

	if _, err := NewClient(srv.URL).Eth().GetFeeHistory(context.Background(), 1, "", nil); err != nil {
		t.Fatalf("GetFeeHistory() with empty block = %v", err)
	}
}

func TestGetPendingTransactions(t *testing.T) {
	block := map[string]interface{}{
		"number": "pending",
		"hash":   "0x0",
		"transactions": []interface{}{
			map[string]interface{}{"hash": "0xpending1", "from": "0xsender", "to": "0xrecipient", "value": "0x0", "gas": "0x5208", "gasPrice": "0x4a817c800", "nonce": "0x1", "input": "0x"},
		},
	}
	srv := newMockRPCServer(block)
	defer srv.Close()

	txs, err := NewClient(srv.URL).Eth().GetPendingTransactions(context.Background())
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
	srv := newMockRPCServer(map[string]interface{}{"number": "pending", "hash": "0x0", "transactions": []interface{}{}})
	defer srv.Close()

	count, err := NewClient(srv.URL).Eth().GetPendingTransactionCount(context.Background())
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
			map[string]interface{}{"hash": "0xtx1", "from": addr, "to": "0xother", "value": "0x0"},
			map[string]interface{}{"hash": "0xtx2", "from": "0xother", "to": "0xunrelated", "value": "0x0"},
		},
	}
	srv := newMockRPCServer(block)
	defer srv.Close()

	txs, err := NewClient(srv.URL).Eth().GetAccountPendingTransactions(context.Background(), addr)
	if err != nil {
		t.Fatalf("GetAccountPendingTransactions() error = %v", err)
	}
	if len(txs) != 1 {
		t.Errorf("GetAccountPendingTransactions() returned %d txs, want 1", len(txs))
	}
}

func TestIsPendingTransaction(t *testing.T) {
	block := map[string]interface{}{
		"number":       "pending",
		"hash":         "0x0",
		"transactions": []interface{}{map[string]interface{}{"hash": "0xtarget", "from": "0xa", "to": "0xb", "value": "0x0"}},
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
