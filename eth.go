package web3

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
)

type Eth struct {
	client *Client
}

func (c *Client) Eth() *Eth {
	return &Eth{client: c}
}

// parseHexBigInt unmarshals a JSON-encoded hex string (e.g. "0x1a") into a *big.Int.
func parseHexBigInt(result json.RawMessage, field string) (*big.Int, error) {
	var hexValue string
	if err := json.Unmarshal(result, &hexValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", field, err)
	}
	val := new(big.Int)
	val.SetString(hexValue[2:], 16)
	return val, nil
}

// parseHexUint64 unmarshals a JSON-encoded hex string into a uint64.
func parseHexUint64(result json.RawMessage, field string) (uint64, error) {
	val, err := parseHexBigInt(result, field)
	if err != nil {
		return 0, err
	}
	return val.Uint64(), nil
}

// parseJSONString unmarshals a JSON-encoded string value.
func parseJSONString(result json.RawMessage, field string) (string, error) {
	var val string
	if err := json.Unmarshal(result, &val); err != nil {
		return "", fmt.Errorf("failed to unmarshal %s: %w", field, err)
	}
	return val, nil
}

func (e *Eth) GetBalance(ctx context.Context, address string, blockNumber BlockParameter) (*big.Int, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	result, err := e.client.Call(ctx, EthGetBalance.String(), []interface{}{address, blockNumber.String()})
	if err != nil {
		return nil, err
	}
	return parseHexBigInt(result, "balance")
}

func (e *Eth) GetBlockNumber(ctx context.Context) (uint64, error) {
	result, err := e.client.Call(ctx, EthGetBlockNumber.String(), []interface{}{})
	if err != nil {
		return 0, err
	}
	return parseHexUint64(result, "block number")
}

func (e *Eth) GetGasPrice(ctx context.Context) (*big.Int, error) {
	result, err := e.client.Call(ctx, EthGetGasPrice.String(), []interface{}{})
	if err != nil {
		return nil, err
	}
	return parseHexBigInt(result, "gas price")
}

func (e *Eth) GetTransactionCount(ctx context.Context, address string, blockNumber BlockParameter) (uint64, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	result, err := e.client.Call(ctx, EthGetTransactionCount.String(), []interface{}{address, blockNumber.String()})
	if err != nil {
		return 0, err
	}
	return parseHexUint64(result, "transaction count")
}

type Block struct {
	Number           string        `json:"number"`
	Hash             string        `json:"hash"`
	ParentHash       string        `json:"parentHash"`
	Nonce            string        `json:"nonce"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	LogsBloom        string        `json:"logsBloom"`
	TransactionsRoot string        `json:"transactionsRoot"`
	StateRoot        string        `json:"stateRoot"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Miner            string        `json:"miner"`
	Difficulty       string        `json:"difficulty"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	ExtraData        string        `json:"extraData"`
	Size             string        `json:"size"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Timestamp        string        `json:"timestamp"`
	Transactions     []interface{} `json:"transactions"`
	Uncles           []string      `json:"uncles"`
}

func (e *Eth) GetBlockByNumber(ctx context.Context, blockNumber BlockParameter, fullTransactions bool) (*Block, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	result, err := e.client.Call(ctx, EthGetBlockByNumber.String(), []interface{}{blockNumber.String(), fullTransactions})
	if err != nil {
		return nil, err
	}
	var block Block
	if err := json.Unmarshal(result, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %w", err)
	}
	return &block, nil
}

func (e *Eth) GetBlockByHash(ctx context.Context, blockHash string, fullTransactions bool) (*Block, error) {
	result, err := e.client.Call(ctx, EthGetBlockByHash.String(), []interface{}{blockHash, fullTransactions})
	if err != nil {
		return nil, err
	}
	var block Block
	if err := json.Unmarshal(result, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %w", err)
	}
	return &block, nil
}

type Transaction struct {
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Input            string `json:"input"`
}

func (e *Eth) GetTransactionByHash(ctx context.Context, txHash string) (*Transaction, error) {
	result, err := e.client.Call(ctx, EthGetTransactionByHash.String(), []interface{}{txHash})
	if err != nil {
		return nil, err
	}
	var tx Transaction
	if err := json.Unmarshal(result, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}
	return &tx, nil
}

type TransactionReceipt struct {
	TransactionHash   string `json:"transactionHash"`
	TransactionIndex  string `json:"transactionIndex"`
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	From              string `json:"from"`
	To                string `json:"to"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	ContractAddress   string `json:"contractAddress"`
	Status            string `json:"status"`
}

func (e *Eth) GetTransactionReceipt(ctx context.Context, txHash string) (*TransactionReceipt, error) {
	result, err := e.client.Call(ctx, EthGetTransactionReceipt.String(), []interface{}{txHash})
	if err != nil {
		return nil, err
	}
	var receipt TransactionReceipt
	if err := json.Unmarshal(result, &receipt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction receipt: %w", err)
	}
	return &receipt, nil
}

func (e *Eth) SendRawTransaction(ctx context.Context, signedTx string) (string, error) {
	result, err := e.client.Call(ctx, EthSendRawTransaction.String(), []interface{}{signedTx})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "transaction hash")
}

func (e *Eth) EstimateGas(ctx context.Context, tx map[string]interface{}) (uint64, error) {
	result, err := e.client.Call(ctx, EthEstimateGas.String(), []interface{}{tx})
	if err != nil {
		return 0, err
	}
	return parseHexUint64(result, "gas estimate")
}

func (e *Eth) Call(ctx context.Context, callObj map[string]interface{}, blockNumber BlockParameter) (string, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	result, err := e.client.Call(ctx, EthCall.String(), []interface{}{callObj, blockNumber.String()})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "call result")
}

// GetPendingTransactions returns pending transactions from the mempool.
func (e *Eth) GetPendingTransactions(ctx context.Context) ([]*Transaction, error) {
	block, err := e.GetBlockByNumber(ctx, BlockPending, true)
	if err != nil {
		return nil, err
	}
	return txsFromBlock(block), nil
}

// txsFromBlock converts the untyped transaction slice from a block into []*Transaction.
func txsFromBlock(block *Block) []*Transaction {
	var txs []*Transaction
	for _, raw := range block.Transactions {
		txData, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		tx := &Transaction{}
		strField := func(key string) string {
			v, _ := txData[key].(string)
			return v
		}
		tx.Hash = strField("hash")
		tx.Nonce = strField("nonce")
		tx.BlockHash = strField("blockHash")
		tx.BlockNumber = strField("blockNumber")
		tx.TransactionIndex = strField("transactionIndex")
		tx.From = strField("from")
		tx.To = strField("to")
		tx.Value = strField("value")
		tx.Gas = strField("gas")
		tx.GasPrice = strField("gasPrice")
		tx.Input = strField("input")
		txs = append(txs, tx)
	}
	return txs
}

// GetPendingTransactionCount returns the number of pending transactions.
func (e *Eth) GetPendingTransactionCount(ctx context.Context) (int, error) {
	txs, err := e.GetPendingTransactions(ctx)
	if err != nil {
		return 0, err
	}
	return len(txs), nil
}

// GetAccountPendingTransactions returns pending transactions for a specific account.
func (e *Eth) GetAccountPendingTransactions(ctx context.Context, address string) ([]*Transaction, error) {
	all, err := e.GetPendingTransactions(ctx)
	if err != nil {
		return nil, err
	}
	var out []*Transaction
	for _, tx := range all {
		if tx.From == address || tx.To == address {
			out = append(out, tx)
		}
	}
	return out, nil
}

// IsPendingTransaction reports whether the given hash is in the pending pool.
func (e *Eth) IsPendingTransaction(ctx context.Context, txHash string) (bool, error) {
	txs, err := e.GetPendingTransactions(ctx)
	if err != nil {
		return false, err
	}
	for _, tx := range txs {
		if tx.Hash == txHash {
			return true, nil
		}
	}
	return false, nil
}

// LogFilter defines filter parameters for eth_getLogs.
type LogFilter struct {
	FromBlock BlockParameter `json:"fromBlock,omitempty"`
	ToBlock   BlockParameter `json:"toBlock,omitempty"`
	Address   interface{}    `json:"address,omitempty"` // string or []string
	Topics    []interface{}  `json:"topics,omitempty"`  // each element: nil, string, or []string
	BlockHash string         `json:"blockHash,omitempty"`
}

// Log represents a single event log entry.
type Log struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	BlockTimestamp   string   `json:"blockTimestamp,omitempty"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

// GetLogs returns event logs matching the given filter.
func (e *Eth) GetLogs(ctx context.Context, filter LogFilter) ([]*Log, error) {
	result, err := e.client.Call(ctx, EthGetLogs.String(), []interface{}{filter})
	if err != nil {
		return nil, err
	}
	var logs []*Log
	if err := json.Unmarshal(result, &logs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal logs: %w", err)
	}
	return logs, nil
}

// GetStorageAt returns the value of a storage slot at a given address.
func (e *Eth) GetStorageAt(ctx context.Context, address, position string, blockNumber BlockParameter) (string, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	result, err := e.client.Call(ctx, EthGetStorageAt.String(), []interface{}{address, position, blockNumber.String()})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "storage value")
}

// GetCode returns the compiled bytecode of a smart contract.
func (e *Eth) GetCode(ctx context.Context, address string, blockNumber BlockParameter) (string, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	result, err := e.client.Call(ctx, EthGetCode.String(), []interface{}{address, blockNumber.String()})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "code")
}

// GetNetVersion returns the current network ID.
func (e *Eth) GetNetVersion(ctx context.Context) (string, error) {
	result, err := e.client.Call(ctx, NetVersion.String(), []interface{}{})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "net version")
}

// GetClientVersion returns the current client version string.
func (e *Eth) GetClientVersion(ctx context.Context) (string, error) {
	result, err := e.client.Call(ctx, Web3ClientVersion.String(), []interface{}{})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "client version")
}

// GetChainID returns the chain ID of the connected network.
func (e *Eth) GetChainID(ctx context.Context) (*big.Int, error) {
	result, err := e.client.Call(ctx, EthChainId.String(), []interface{}{})
	if err != nil {
		return nil, err
	}
	return parseHexBigInt(result, "chain ID")
}

// GetMaxPriorityFeePerGas returns the current maxPriorityFeePerGas (EIP-1559 tip).
func (e *Eth) GetMaxPriorityFeePerGas(ctx context.Context) (*big.Int, error) {
	result, err := e.client.Call(ctx, EthMaxPriorityFeePerGas.String(), []interface{}{})
	if err != nil {
		return nil, err
	}
	return parseHexBigInt(result, "max priority fee per gas")
}

// FeeHistoryResult holds the response from eth_feeHistory.
type FeeHistoryResult struct {
	OldestBlock       string     `json:"oldestBlock"`
	BaseFeePerGas     []string   `json:"baseFeePerGas"`
	GasUsedRatio      []float64  `json:"gasUsedRatio"`
	Reward            [][]string `json:"reward,omitempty"`
	BaseFeePerBlobGas []string   `json:"baseFeePerBlobGas,omitempty"` // EIP-4844
	BlobGasUsedRatio  []float64  `json:"blobGasUsedRatio,omitempty"`  // EIP-4844
}

// GetFeeHistory returns historical gas fee data for the given block range.
func (e *Eth) GetFeeHistory(ctx context.Context, blockCount uint64, newestBlock BlockParameter, rewardPercentiles []float64) (*FeeHistoryResult, error) {
	if newestBlock == "" {
		newestBlock = BlockLatest
	}
	percentiles := rewardPercentiles
	if len(percentiles) == 0 {
		percentiles = []float64{}
	}
	result, err := e.client.Call(ctx, EthFeeHistory.String(), []interface{}{ToHex(blockCount), newestBlock.String(), percentiles})
	if err != nil {
		return nil, err
	}
	var feeHistory FeeHistoryResult
	if err := json.Unmarshal(result, &feeHistory); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fee history: %w", err)
	}
	return &feeHistory, nil
}
