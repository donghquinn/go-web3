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