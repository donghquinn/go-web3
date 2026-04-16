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

func (e *Eth) GetBalance(ctx context.Context, address string, blockNumber BlockParameter) (*big.Int, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	
	result, err := e.client.Call(ctx, EthGetBalance.String(), []interface{}{address, blockNumber.String()})
	if err != nil {
		return nil, err
	}

	var hexValue string
	if err := json.Unmarshal(result, &hexValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal balance: %w", err)
	}

	balance := new(big.Int)
	balance.SetString(hexValue[2:], 16)
	return balance, nil
}

func (e *Eth) GetBlockNumber(ctx context.Context) (uint64, error) {
	result, err := e.client.Call(ctx, EthGetBlockNumber.String(), []interface{}{})
	if err != nil {
		return 0, err
	}

	var hexValue string
	if err := json.Unmarshal(result, &hexValue); err != nil {
		return 0, fmt.Errorf("failed to unmarshal block number: %w", err)
	}

	blockNumber := new(big.Int)
	blockNumber.SetString(hexValue[2:], 16)
	return blockNumber.Uint64(), nil
}

func (e *Eth) GetGasPrice(ctx context.Context) (*big.Int, error) {
	result, err := e.client.Call(ctx, EthGetGasPrice.String(), []interface{}{})
	if err != nil {
		return nil, err
	}

	var hexValue string
	if err := json.Unmarshal(result, &hexValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gas price: %w", err)
	}

	gasPrice := new(big.Int)
	gasPrice.SetString(hexValue[2:], 16)
	return gasPrice, nil
}

func (e *Eth) GetTransactionCount(ctx context.Context, address string, blockNumber BlockParameter) (uint64, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	
	result, err := e.client.Call(ctx, EthGetTransactionCount.String(), []interface{}{address, blockNumber.String()})
	if err != nil {
		return 0, err
	}

	var hexValue string
	if err := json.Unmarshal(result, &hexValue); err != nil {
		return 0, fmt.Errorf("failed to unmarshal transaction count: %w", err)
	}

	nonce := new(big.Int)
	nonce.SetString(hexValue[2:], 16)
	return nonce.Uint64(), nil
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

	var txHash string
	if err := json.Unmarshal(result, &txHash); err != nil {
		return "", fmt.Errorf("failed to unmarshal transaction hash: %w", err)
	}

	return txHash, nil
}

func (e *Eth) EstimateGas(ctx context.Context, tx map[string]interface{}) (uint64, error) {
	result, err := e.client.Call(ctx, EthEstimateGas.String(), []interface{}{tx})
	if err != nil {
		return 0, err
	}

	var hexValue string
	if err := json.Unmarshal(result, &hexValue); err != nil {
		return 0, fmt.Errorf("failed to unmarshal gas estimate: %w", err)
	}

	gasEstimate := new(big.Int)
	gasEstimate.SetString(hexValue[2:], 16)
	return gasEstimate.Uint64(), nil
}

func (e *Eth) Call(ctx context.Context, callObj map[string]interface{}, blockNumber BlockParameter) (string, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	
	result, err := e.client.Call(ctx, EthCall.String(), []interface{}{callObj, blockNumber.String()})
	if err != nil {
		return "", err
	}

	var data string
	if err := json.Unmarshal(result, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal call result: %w", err)
	}

	return data, nil
}

// GetPendingTransactions returns pending transactions from the mempool
func (e *Eth) GetPendingTransactions(ctx context.Context) ([]*Transaction, error) {
	// Get the pending block with full transaction details
	block, err := e.GetBlockByNumber(ctx, BlockPending, true)
	if err != nil {
		return nil, err
	}
	
	// Convert interface{} transactions to Transaction structs
	var pendingTxs []*Transaction
	for _, txInterface := range block.Transactions {
		if txData, ok := txInterface.(map[string]interface{}); ok {
			tx := &Transaction{}
			
			// Parse transaction fields with proper error handling
			if hash, ok := txData["hash"].(string); ok {
				tx.Hash = hash
			}
			if nonce, ok := txData["nonce"].(string); ok {
				tx.Nonce = nonce
			}
			if blockHash, ok := txData["blockHash"].(string); ok {
				tx.BlockHash = blockHash
			}
			if blockNumber, ok := txData["blockNumber"].(string); ok {
				tx.BlockNumber = blockNumber
			}
			if transactionIndex, ok := txData["transactionIndex"].(string); ok {
				tx.TransactionIndex = transactionIndex
			}
			if from, ok := txData["from"].(string); ok {
				tx.From = from
			}
			if to, ok := txData["to"].(string); ok {
				tx.To = to
			}
			if value, ok := txData["value"].(string); ok {
				tx.Value = value
			}
			if gas, ok := txData["gas"].(string); ok {
				tx.Gas = gas
			}
			if gasPrice, ok := txData["gasPrice"].(string); ok {
				tx.GasPrice = gasPrice
			}
			if input, ok := txData["input"].(string); ok {
				tx.Input = input
			}
			
			pendingTxs = append(pendingTxs, tx)
		}
	}
	
	return pendingTxs, nil
}

// GetPendingTransactionCount returns the number of pending transactions
func (e *Eth) GetPendingTransactionCount(ctx context.Context) (int, error) {
	pendingTxs, err := e.GetPendingTransactions(ctx)
	if err != nil {
		return 0, err
	}
	return len(pendingTxs), nil
}

// GetAccountPendingTransactions returns pending transactions for a specific account
func (e *Eth) GetAccountPendingTransactions(ctx context.Context, address string) ([]*Transaction, error) {
	allPendingTxs, err := e.GetPendingTransactions(ctx)
	if err != nil {
		return nil, err
	}
	
	var accountTxs []*Transaction
	for _, tx := range allPendingTxs {
		if tx.From == address || tx.To == address {
			accountTxs = append(accountTxs, tx)
		}
	}
	
	return accountTxs, nil
}

// IsPendingTransaction checks if a transaction hash is in the pending pool
func (e *Eth) IsPendingTransaction(ctx context.Context, txHash string) (bool, error) {
	pendingTxs, err := e.GetPendingTransactions(ctx)
	if err != nil {
		return false, err
	}

	for _, tx := range pendingTxs {
		if tx.Hash == txHash {
			return true, nil
		}
	}

	return false, nil
}

// LogFilter defines filter parameters for eth_getLogs
type LogFilter struct {
	FromBlock BlockParameter `json:"fromBlock,omitempty"`
	ToBlock   BlockParameter `json:"toBlock,omitempty"`
	Address   interface{}    `json:"address,omitempty"` // string or []string
	Topics    []interface{}  `json:"topics,omitempty"`  // each element: nil, string, or []string
	BlockHash string         `json:"blockHash,omitempty"`
}

// Log represents a single event log entry
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

// GetLogs returns event logs matching the given filter
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

// GetStorageAt returns the value of a storage slot at a given address
func (e *Eth) GetStorageAt(ctx context.Context, address string, position string, blockNumber BlockParameter) (string, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}

	result, err := e.client.Call(ctx, EthGetStorageAt.String(), []interface{}{address, position, blockNumber.String()})
	if err != nil {
		return "", err
	}

	var value string
	if err := json.Unmarshal(result, &value); err != nil {
		return "", fmt.Errorf("failed to unmarshal storage value: %w", err)
	}

	return value, nil
}

// GetCode returns the compiled bytecode of a smart contract
func (e *Eth) GetCode(ctx context.Context, address string, blockNumber BlockParameter) (string, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}

	result, err := e.client.Call(ctx, EthGetCode.String(), []interface{}{address, blockNumber.String()})
	if err != nil {
		return "", err
	}

	var code string
	if err := json.Unmarshal(result, &code); err != nil {
		return "", fmt.Errorf("failed to unmarshal code: %w", err)
	}

	return code, nil
}

// GetNetVersion returns the current network ID
func (e *Eth) GetNetVersion(ctx context.Context) (string, error) {
	result, err := e.client.Call(ctx, NetVersion.String(), []interface{}{})
	if err != nil {
		return "", err
	}

	var version string
	if err := json.Unmarshal(result, &version); err != nil {
		return "", fmt.Errorf("failed to unmarshal net version: %w", err)
	}

	return version, nil
}

// GetClientVersion returns the current client version string
func (e *Eth) GetClientVersion(ctx context.Context) (string, error) {
	result, err := e.client.Call(ctx, Web3ClientVersion.String(), []interface{}{})
	if err != nil {
		return "", err
	}

	var version string
	if err := json.Unmarshal(result, &version); err != nil {
		return "", fmt.Errorf("failed to unmarshal client version: %w", err)
	}

	return version, nil
}

// GetChainID returns the chain ID of the connected network
func (e *Eth) GetChainID(ctx context.Context) (*big.Int, error) {
	result, err := e.client.Call(ctx, EthChainId.String(), []interface{}{})
	if err != nil {
		return nil, err
	}

	var hexValue string
	if err := json.Unmarshal(result, &hexValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chain ID: %w", err)
	}

	chainID := new(big.Int)
	chainID.SetString(hexValue[2:], 16)
	return chainID, nil
}

// GetMaxPriorityFeePerGas returns the current maxPriorityFeePerGas (EIP-1559 tip)
func (e *Eth) GetMaxPriorityFeePerGas(ctx context.Context) (*big.Int, error) {
	result, err := e.client.Call(ctx, EthMaxPriorityFeePerGas.String(), []interface{}{})
	if err != nil {
		return nil, err
	}

	var hexValue string
	if err := json.Unmarshal(result, &hexValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal max priority fee per gas: %w", err)
	}

	fee := new(big.Int)
	fee.SetString(hexValue[2:], 16)
	return fee, nil
}

// FeeHistoryResult holds the response from eth_feeHistory
type FeeHistoryResult struct {
	OldestBlock        string     `json:"oldestBlock"`
	BaseFeePerGas      []string   `json:"baseFeePerGas"`
	GasUsedRatio       []float64  `json:"gasUsedRatio"`
	Reward             [][]string `json:"reward,omitempty"`
	BaseFeePerBlobGas  []string   `json:"baseFeePerBlobGas,omitempty"`  // EIP-4844
	BlobGasUsedRatio   []float64  `json:"blobGasUsedRatio,omitempty"`   // EIP-4844
}

// GetFeeHistory returns historical gas fee data for the given block range.
// blockCount is the number of blocks to include (hex string or uint64 as string).
// rewardPercentiles is an optional list of percentiles (0-100) for priority fee sampling.
func (e *Eth) GetFeeHistory(ctx context.Context, blockCount uint64, newestBlock BlockParameter, rewardPercentiles []float64) (*FeeHistoryResult, error) {
	if newestBlock == "" {
		newestBlock = BlockLatest
	}

	params := []interface{}{ToHex(blockCount), newestBlock.String()}
	if len(rewardPercentiles) > 0 {
		params = append(params, rewardPercentiles)
	} else {
		params = append(params, []float64{})
	}

	result, err := e.client.Call(ctx, EthFeeHistory.String(), params)
	if err != nil {
		return nil, err
	}

	var feeHistory FeeHistoryResult
	if err := json.Unmarshal(result, &feeHistory); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fee history: %w", err)
	}

	return &feeHistory, nil
}