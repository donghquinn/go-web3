package web3

import (
	"context"
	"reflect"
	"testing"
)

// ethCall invokes a single Eth method and discards its typed result.
type ethCall func(ctx context.Context, e *Eth) (interface{}, error)

// extendedCalls lists every API added for web3.js parity, with the RPC method and params each should send.
var extendedCalls = []struct {
	name   string
	method RPCMethod
	params []interface{}
	call   ethCall
}{
	{"GetBlockTransactionCountByNumber", EthGetBlockTransactionCountByNumber, []interface{}{"0x10"},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.GetBlockTransactionCountByNumber(ctx, BlockNumber(16))
		}},
	{"GetBlockTransactionCountByHash", EthGetBlockTransactionCountByHash, []interface{}{"0xblock"},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.GetBlockTransactionCountByHash(ctx, "0xblock")
		}},
	{"GetUncleCountByBlockNumber", EthGetUncleCountByBlockNumber, []interface{}{"latest"},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.GetUncleCountByBlockNumber(ctx, "")
		}},
	{"GetUncleCountByBlockHash", EthGetUncleCountByBlockHash, []interface{}{"0xblock"},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.GetUncleCountByBlockHash(ctx, "0xblock")
		}},
	{"GetUncleByBlockNumberAndIndex", EthGetUncleByBlockNumberAndIndex, []interface{}{"latest", "0x1"},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.GetUncleByBlockNumberAndIndex(ctx, "", 1)
		}},
	{"GetUncleByBlockHashAndIndex", EthGetUncleByBlockHashAndIndex, []interface{}{"0xblock", "0x0"},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.GetUncleByBlockHashAndIndex(ctx, "0xblock", 0)
		}},
	{"GetTransactionByBlockNumberAndIndex", EthGetTransactionByBlockNumberAndIndex, []interface{}{"latest", "0x2"},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.GetTransactionByBlockNumberAndIndex(ctx, "", 2)
		}},
	{"GetTransactionByBlockHashAndIndex", EthGetTransactionByBlockHashAndIndex, []interface{}{"0xblock", "0x3"},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.GetTransactionByBlockHashAndIndex(ctx, "0xblock", 3)
		}},
	{"GetBlockReceipts", EthGetBlockReceipts, []interface{}{"latest"},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.GetBlockReceipts(ctx, "")
		}},
	{"IsSyncing", EthSyncing, []interface{}{},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.IsSyncing(ctx) }},
	{"GetCoinbase", EthCoinbase, []interface{}{},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.GetCoinbase(ctx) }},
	{"IsMining", EthMining, []interface{}{},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.IsMining(ctx) }},
	{"GetHashRate", EthHashrate, []interface{}{},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.GetHashRate(ctx) }},
	{"GetAccounts", EthAccounts, []interface{}{},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.GetAccounts(ctx) }},
	{"GetProtocolVersion", EthProtocolVersion, []interface{}{},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.GetProtocolVersion(ctx) }},
	{"GetBlobBaseFee", EthBlobBaseFee, []interface{}{},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.GetBlobBaseFee(ctx) }},
	{"IsListening", NetListening, []interface{}{},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.IsListening(ctx) }},
	{"GetPeerCount", NetPeerCount, []interface{}{},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.GetPeerCount(ctx) }},
	{"Sha3", Web3Sha3, []interface{}{"0x68656c6c6f"},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.Sha3(ctx, "0x68656c6c6f") }},
	{"GetProof", EthGetProof, []interface{}{"0xaddr", []interface{}{}, "latest"},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.GetProof(ctx, "0xaddr", nil, "") }},
	{"CreateAccessList", EthCreateAccessList, []interface{}{map[string]interface{}{"to": "0xc"}, "latest"},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.CreateAccessList(ctx, map[string]interface{}{"to": "0xc"}, "")
		}},
	{"Sign", EthSign, []interface{}{"0xaddr", "0xdead"},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.Sign(ctx, "0xaddr", "0xdead") }},
	{"SignTransaction", EthSignTransaction, []interface{}{map[string]interface{}{"from": "0xa"}},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.SignTransaction(ctx, map[string]interface{}{"from": "0xa"})
		}},
	{"SendTransaction", EthSendTransaction, []interface{}{map[string]interface{}{"from": "0xa"}},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.SendTransaction(ctx, map[string]interface{}{"from": "0xa"})
		}},
	{"NewFilter", EthNewFilter, []interface{}{map[string]interface{}{"fromBlock": "latest", "address": "0xc"}},
		func(ctx context.Context, e *Eth) (interface{}, error) {
			return e.NewFilter(ctx, LogFilter{FromBlock: BlockLatest, Address: "0xc"})
		}},
	{"NewBlockFilter", EthNewBlockFilter, []interface{}{},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.NewBlockFilter(ctx) }},
	{"NewPendingTransactionFilter", EthNewPendingTransactionFilter, []interface{}{},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.NewPendingTransactionFilter(ctx) }},
	{"GetFilterLogChanges", EthGetFilterChanges, []interface{}{"0x1"},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.GetFilterLogChanges(ctx, "0x1") }},
	{"GetFilterHashChanges", EthGetFilterChanges, []interface{}{"0x1"},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.GetFilterHashChanges(ctx, "0x1") }},
	{"GetFilterLogs", EthGetFilterLogs, []interface{}{"0x1"},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.GetFilterLogs(ctx, "0x1") }},
	{"UninstallFilter", EthUninstallFilter, []interface{}{"0x1"},
		func(ctx context.Context, e *Eth) (interface{}, error) { return e.UninstallFilter(ctx, "0x1") }},
}

func TestExtendedCallsRequest(t *testing.T) {
	for _, tc := range extendedCalls {
		t.Run(tc.name, func(t *testing.T) {
			var got RPCRequest
			srv := newCaptureServer(nil, &got)
			defer srv.Close()

			tc.call(context.Background(), NewClient(srv.URL).Eth())

			if got.Method != tc.method.String() {
				t.Errorf("method = %q, want %q", got.Method, tc.method)
			}
			if !reflect.DeepEqual(got.Params, tc.params) {
				t.Errorf("params = %#v, want %#v", got.Params, tc.params)
			}
		})
	}
}

func TestExtendedCallsRPCError(t *testing.T) {
	srv := newMockRPCErrorServer(-32601, "method not found")
	defer srv.Close()
	eth := NewClient(srv.URL).Eth()

	for _, tc := range extendedCalls {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.call(context.Background(), eth); err == nil {
				t.Fatal("expected RPC error")
			}
		})
	}
}

func TestExtendedCallsUnmarshalError(t *testing.T) {
	// An object is not a valid result for any of these methods except eth_syncing, which is tested separately.
	srv := newMockRPCServer(map[string]interface{}{"accessList": 1, "raw": 1, "address": 1, "number": 1, "hash": 1})
	defer srv.Close()
	eth := NewClient(srv.URL).Eth()

	for _, tc := range extendedCalls {
		if tc.method == EthSyncing {
			continue
		}
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.call(context.Background(), eth); err == nil {
				t.Fatal("expected unmarshal error")
			}
		})
	}
}

func TestGetBlockTransactionCount(t *testing.T) {
	srv := newMockRPCServer("0x2a")
	defer srv.Close()
	eth := NewClient(srv.URL).Eth()

	byNum, err := eth.GetBlockTransactionCountByNumber(context.Background(), BlockLatest)
	if err != nil || byNum != 42 {
		t.Errorf("GetBlockTransactionCountByNumber() = %d, %v; want 42", byNum, err)
	}
	byHash, err := eth.GetBlockTransactionCountByHash(context.Background(), "0xblock")
	if err != nil || byHash != 42 {
		t.Errorf("GetBlockTransactionCountByHash() = %d, %v; want 42", byHash, err)
	}
}

func TestGetUncleCount(t *testing.T) {
	srv := newMockRPCServer("0x1")
	defer srv.Close()
	eth := NewClient(srv.URL).Eth()

	byNum, err := eth.GetUncleCountByBlockNumber(context.Background(), BlockLatest)
	if err != nil || byNum != 1 {
		t.Errorf("GetUncleCountByBlockNumber() = %d, %v; want 1", byNum, err)
	}
	byHash, err := eth.GetUncleCountByBlockHash(context.Background(), "0xblock")
	if err != nil || byHash != 1 {
		t.Errorf("GetUncleCountByBlockHash() = %d, %v; want 1", byHash, err)
	}
}

func TestGetUncle(t *testing.T) {
	srv := newMockRPCServer(map[string]interface{}{"hash": "0xuncle", "number": "0x10"})
	defer srv.Close()
	eth := NewClient(srv.URL).Eth()

	uncle, err := eth.GetUncleByBlockNumberAndIndex(context.Background(), BlockLatest, 0)
	if err != nil || uncle == nil || uncle.Hash != "0xuncle" {
		t.Errorf("GetUncleByBlockNumberAndIndex() = %+v, %v", uncle, err)
	}
	uncle, err = eth.GetUncleByBlockHashAndIndex(context.Background(), "0xblock", 0)
	if err != nil || uncle == nil || uncle.Hash != "0xuncle" {
		t.Errorf("GetUncleByBlockHashAndIndex() = %+v, %v", uncle, err)
	}
}

func TestGetTransactionFromBlock(t *testing.T) {
	srv := newMockRPCServer(map[string]interface{}{
		"hash":                 "0xtx",
		"type":                 "0x2",
		"maxFeePerGas":         "0x77359400",
		"maxPriorityFeePerGas": "0x3b9aca00",
		"accessList":           []interface{}{map[string]interface{}{"address": "0xc", "storageKeys": []string{"0x0"}}},
	})
	defer srv.Close()
	eth := NewClient(srv.URL).Eth()

	tx, err := eth.GetTransactionByBlockNumberAndIndex(context.Background(), BlockLatest, 0)
	if err != nil {
		t.Fatalf("GetTransactionByBlockNumberAndIndex() error = %v", err)
	}
	if tx.Hash != "0xtx" || tx.Type != "0x2" || tx.MaxFeePerGas != "0x77359400" || len(tx.AccessList) != 1 {
		t.Errorf("GetTransactionByBlockNumberAndIndex() = %+v", tx)
	}
	tx, err = eth.GetTransactionByBlockHashAndIndex(context.Background(), "0xblock", 0)
	if err != nil || tx == nil || tx.Hash != "0xtx" {
		t.Errorf("GetTransactionByBlockHashAndIndex() = %+v, %v", tx, err)
	}
}

func TestNullResultReturnsNil(t *testing.T) {
	srv := newMockRPCServer(nil)
	defer srv.Close()
	eth := NewClient(srv.URL).Eth()
	ctx := context.Background()

	if tx, err := eth.GetTransactionByBlockHashAndIndex(ctx, "0xblock", 9); err != nil || tx != nil {
		t.Errorf("GetTransactionByBlockHashAndIndex() on null = %+v, %v; want nil, nil", tx, err)
	}
	if uncle, err := eth.GetUncleByBlockHashAndIndex(ctx, "0xblock", 9); err != nil || uncle != nil {
		t.Errorf("GetUncleByBlockHashAndIndex() on null = %+v, %v; want nil, nil", uncle, err)
	}
	if receipts, err := eth.GetBlockReceipts(ctx, "0xblock"); err != nil || receipts != nil {
		t.Errorf("GetBlockReceipts() on null = %+v, %v; want nil, nil", receipts, err)
	}
}

func TestGetBlockReceipts(t *testing.T) {
	srv := newMockRPCServer([]interface{}{
		map[string]interface{}{
			"transactionHash":   "0xtx",
			"status":            "0x1",
			"effectiveGasPrice": "0x3b9aca00",
			"logs":              []interface{}{map[string]interface{}{"address": "0xc", "topics": []string{"0xt"}}},
		},
	})
	defer srv.Close()

	receipts, err := NewClient(srv.URL).Eth().GetBlockReceipts(context.Background(), BlockNumber(1))
	if err != nil {
		t.Fatalf("GetBlockReceipts() error = %v", err)
	}
	if len(receipts) != 1 || receipts[0].EffectiveGasPrice != "0x3b9aca00" || len(receipts[0].Logs) != 1 {
		t.Errorf("GetBlockReceipts() = %+v", receipts)
	}
}

func TestIsSyncing(t *testing.T) {
	notSyncing := newMockRPCServer(false)
	defer notSyncing.Close()
	status, err := NewClient(notSyncing.URL).Eth().IsSyncing(context.Background())
	if err != nil || status != nil {
		t.Errorf("IsSyncing() when synced = %+v, %v; want nil, nil", status, err)
	}

	syncing := newMockRPCServer(map[string]interface{}{"startingBlock": "0x0", "currentBlock": "0x10", "highestBlock": "0x20"})
	defer syncing.Close()
	status, err = NewClient(syncing.URL).Eth().IsSyncing(context.Background())
	if err != nil || status == nil || status.CurrentBlock != "0x10" || status.HighestBlock != "0x20" {
		t.Errorf("IsSyncing() when syncing = %+v, %v", status, err)
	}

	bad := newMockRPCServer("syncing")
	defer bad.Close()
	if _, err := NewClient(bad.URL).Eth().IsSyncing(context.Background()); err == nil {
		t.Error("IsSyncing() expected unmarshal error")
	}
}

func TestNodeInfo(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		EthCoinbase.String():        "0xcoinbase",
		EthMining.String():          true,
		EthHashrate.String():        "0x38a",
		EthAccounts.String():        []string{"0xa", "0xb"},
		EthProtocolVersion.String(): "0x41",
		EthBlobBaseFee.String():     "0x1",
		NetListening.String():       true,
		NetPeerCount.String():       "0x19",
		Web3Sha3.String():           "0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8",
	})
	defer srv.Close()
	eth := NewClient(srv.URL).Eth()
	ctx := context.Background()

	if v, err := eth.GetCoinbase(ctx); err != nil || v != "0xcoinbase" {
		t.Errorf("GetCoinbase() = %q, %v", v, err)
	}
	if v, err := eth.IsMining(ctx); err != nil || !v {
		t.Errorf("IsMining() = %v, %v", v, err)
	}
	if v, err := eth.GetHashRate(ctx); err != nil || v.Int64() != 906 {
		t.Errorf("GetHashRate() = %v, %v", v, err)
	}
	if v, err := eth.GetAccounts(ctx); err != nil || !reflect.DeepEqual(v, []string{"0xa", "0xb"}) {
		t.Errorf("GetAccounts() = %v, %v", v, err)
	}
	if v, err := eth.GetProtocolVersion(ctx); err != nil || v != "0x41" {
		t.Errorf("GetProtocolVersion() = %q, %v", v, err)
	}
	if v, err := eth.GetBlobBaseFee(ctx); err != nil || v.Int64() != 1 {
		t.Errorf("GetBlobBaseFee() = %v, %v", v, err)
	}
	if v, err := eth.IsListening(ctx); err != nil || !v {
		t.Errorf("IsListening() = %v, %v", v, err)
	}
	if v, err := eth.GetPeerCount(ctx); err != nil || v != 25 {
		t.Errorf("GetPeerCount() = %d, %v", v, err)
	}
	if v, err := eth.Sha3(ctx, "0x68656c6c6f20776f726c64"); err != nil || v != "0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8" {
		t.Errorf("Sha3() = %q, %v", v, err)
	}
}

func TestGetProof(t *testing.T) {
	srv := newMockRPCServer(map[string]interface{}{
		"address":      "0xaddr",
		"accountProof": []string{"0xp1", "0xp2"},
		"balance":      "0x0",
		"codeHash":     "0xcode",
		"nonce":        "0x1",
		"storageHash":  "0xstorage",
		"storageProof": []interface{}{map[string]interface{}{"key": "0x0", "value": "0x5", "proof": []string{"0xs"}}},
	})
	defer srv.Close()

	proof, err := NewClient(srv.URL).Eth().GetProof(context.Background(), "0xaddr", []string{"0x0"}, BlockLatest)
	if err != nil {
		t.Fatalf("GetProof() error = %v", err)
	}
	if len(proof.AccountProof) != 2 || len(proof.StorageProof) != 1 || proof.StorageProof[0].Value != "0x5" {
		t.Errorf("GetProof() = %+v", proof)
	}
}

func TestCreateAccessList(t *testing.T) {
	srv := newMockRPCServer(map[string]interface{}{
		"accessList": []interface{}{map[string]interface{}{"address": "0xc", "storageKeys": []string{"0x0", "0x1"}}},
		"gasUsed":    "0x5208",
	})
	defer srv.Close()

	res, err := NewClient(srv.URL).Eth().CreateAccessList(context.Background(), map[string]interface{}{"to": "0xc"}, BlockLatest)
	if err != nil {
		t.Fatalf("CreateAccessList() error = %v", err)
	}
	if res.GasUsed != "0x5208" || len(res.AccessList) != 1 || len(res.AccessList[0].StorageKeys) != 2 {
		t.Errorf("CreateAccessList() = %+v", res)
	}
}

func TestNodeSigning(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		EthSign.String():            "0xsig",
		EthSignTransaction.String(): map[string]interface{}{"raw": "0xraw", "tx": map[string]interface{}{"hash": "0xtx"}},
		EthSendTransaction.String(): "0xtx",
	})
	defer srv.Close()
	eth := NewClient(srv.URL).Eth()
	ctx := context.Background()

	if v, err := eth.Sign(ctx, "0xaddr", "0xdead"); err != nil || v != "0xsig" {
		t.Errorf("Sign() = %q, %v", v, err)
	}
	signed, err := eth.SignTransaction(ctx, map[string]interface{}{"from": "0xa"})
	if err != nil || signed.Raw != "0xraw" || signed.Tx.Hash != "0xtx" {
		t.Errorf("SignTransaction() = %+v, %v", signed, err)
	}
	if v, err := eth.SendTransaction(ctx, map[string]interface{}{"from": "0xa"}); err != nil || v != "0xtx" {
		t.Errorf("SendTransaction() = %q, %v", v, err)
	}
}

func TestFilters(t *testing.T) {
	srv := newDispatchServer(map[string]interface{}{
		EthNewFilter.String():                   "0x1",
		EthNewBlockFilter.String():              "0x2",
		EthNewPendingTransactionFilter.String(): "0x3",
		EthGetFilterLogs.String():               []interface{}{map[string]interface{}{"address": "0xc", "logIndex": "0x0"}},
		EthUninstallFilter.String():             true,
	})
	defer srv.Close()
	eth := NewClient(srv.URL).Eth()
	ctx := context.Background()

	if id, err := eth.NewFilter(ctx, LogFilter{Address: "0xc"}); err != nil || id != "0x1" {
		t.Errorf("NewFilter() = %q, %v", id, err)
	}
	if id, err := eth.NewBlockFilter(ctx); err != nil || id != "0x2" {
		t.Errorf("NewBlockFilter() = %q, %v", id, err)
	}
	if id, err := eth.NewPendingTransactionFilter(ctx); err != nil || id != "0x3" {
		t.Errorf("NewPendingTransactionFilter() = %q, %v", id, err)
	}
	if logs, err := eth.GetFilterLogs(ctx, "0x1"); err != nil || len(logs) != 1 || logs[0].Address != "0xc" {
		t.Errorf("GetFilterLogs() = %+v, %v", logs, err)
	}
	if ok, err := eth.UninstallFilter(ctx, "0x1"); err != nil || !ok {
		t.Errorf("UninstallFilter() = %v, %v", ok, err)
	}
}

func TestGetFilterChanges(t *testing.T) {
	logSrv := newMockRPCServer([]interface{}{map[string]interface{}{"address": "0xc", "data": "0x"}})
	defer logSrv.Close()
	logs, err := NewClient(logSrv.URL).Eth().GetFilterLogChanges(context.Background(), "0x1")
	if err != nil || len(logs) != 1 || logs[0].Address != "0xc" {
		t.Errorf("GetFilterLogChanges() = %+v, %v", logs, err)
	}

	hashSrv := newMockRPCServer([]string{"0xh1", "0xh2"})
	defer hashSrv.Close()
	hashes, err := NewClient(hashSrv.URL).Eth().GetFilterHashChanges(context.Background(), "0x2")
	if err != nil || !reflect.DeepEqual(hashes, []string{"0xh1", "0xh2"}) {
		t.Errorf("GetFilterHashChanges() = %v, %v", hashes, err)
	}
}

func TestGetPendingTransactionsTypedFields(t *testing.T) {
	srv := newMockRPCServer(map[string]interface{}{
		"number": "pending",
		"transactions": []interface{}{
			map[string]interface{}{"hash": "0xtx", "type": "0x2", "maxFeePerGas": "0x1", "chainId": "0x1"},
			map[string]interface{}{"hash": 5}, // malformed entry is skipped
		},
	})
	defer srv.Close()

	txs, err := NewClient(srv.URL).Eth().GetPendingTransactions(context.Background())
	if err != nil {
		t.Fatalf("GetPendingTransactions() error = %v", err)
	}
	if len(txs) != 1 || txs[0].Type != "0x2" || txs[0].MaxFeePerGas != "0x1" || txs[0].ChainID != "0x1" {
		t.Errorf("GetPendingTransactions() = %+v", txs)
	}
}

func TestGetBlockTransactionCountDefaultBlock(t *testing.T) {
	var got RPCRequest
	srv := newCaptureServer("0x0", &got)
	defer srv.Close()

	if _, err := NewClient(srv.URL).Eth().GetBlockTransactionCountByNumber(context.Background(), ""); err != nil {
		t.Fatalf("GetBlockTransactionCountByNumber() error = %v", err)
	}
	if got.Params[0] != "latest" {
		t.Errorf("default block = %v, want latest", got.Params[0])
	}
}

func TestParseHexBigIntInvalid(t *testing.T) {
	for _, in := range []string{`null`, `""`, `"0x"`, `"123"`, `"0xzz"`} {
		if _, err := parseHexBigInt([]byte(in), "value"); err == nil {
			t.Errorf("parseHexBigInt(%s) expected error", in)
		}
	}
}
