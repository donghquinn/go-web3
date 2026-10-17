package web3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// isNullResult reports whether an RPC result is JSON null (e.g. an unknown block or transaction).
func isNullResult(result json.RawMessage) bool {
	return len(result) == 0 || bytes.Equal(bytes.TrimSpace(result), []byte("null"))
}

// unmarshalOptional decodes result into a new T, returning nil when the node answered null.
func unmarshalOptional[T any](result json.RawMessage, field string) (*T, error) {
	if isNullResult(result) {
		return nil, nil
	}
	var v T
	if err := json.Unmarshal(result, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", field, err)
	}
	return &v, nil
}

// GetBlockTransactionCountByNumber returns the number of transactions in the given block.
func (e *Eth) GetBlockTransactionCountByNumber(ctx context.Context, blockNumber BlockParameter) (uint64, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	result, err := e.client.Call(ctx, EthGetBlockTransactionCountByNumber.String(), []interface{}{blockNumber.String()})
	if err != nil {
		return 0, err
	}
	return parseHexUint64(result, "block transaction count")
}

// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
func (e *Eth) GetBlockTransactionCountByHash(ctx context.Context, blockHash string) (uint64, error) {
	result, err := e.client.Call(ctx, EthGetBlockTransactionCountByHash.String(), []interface{}{blockHash})
	if err != nil {
		return 0, err
	}
	return parseHexUint64(result, "block transaction count")
}

// GetUncleCountByBlockNumber returns the number of uncles in the given block.
func (e *Eth) GetUncleCountByBlockNumber(ctx context.Context, blockNumber BlockParameter) (uint64, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	result, err := e.client.Call(ctx, EthGetUncleCountByBlockNumber.String(), []interface{}{blockNumber.String()})
	if err != nil {
		return 0, err
	}
	return parseHexUint64(result, "uncle count")
}

// GetUncleCountByBlockHash returns the number of uncles in the block with the given hash.
func (e *Eth) GetUncleCountByBlockHash(ctx context.Context, blockHash string) (uint64, error) {
	result, err := e.client.Call(ctx, EthGetUncleCountByBlockHash.String(), []interface{}{blockHash})
	if err != nil {
		return 0, err
	}
	return parseHexUint64(result, "uncle count")
}

// GetUncleByBlockNumberAndIndex returns an uncle of the given block, or nil if it does not exist.
func (e *Eth) GetUncleByBlockNumberAndIndex(ctx context.Context, blockNumber BlockParameter, index uint64) (*Block, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	result, err := e.client.Call(ctx, EthGetUncleByBlockNumberAndIndex.String(), []interface{}{blockNumber.String(), ToHex(index)})
	if err != nil {
		return nil, err
	}
	return unmarshalOptional[Block](result, "uncle")
}

// GetUncleByBlockHashAndIndex returns an uncle of the block with the given hash, or nil if it does not exist.
func (e *Eth) GetUncleByBlockHashAndIndex(ctx context.Context, blockHash string, index uint64) (*Block, error) {
	result, err := e.client.Call(ctx, EthGetUncleByBlockHashAndIndex.String(), []interface{}{blockHash, ToHex(index)})
	if err != nil {
		return nil, err
	}
	return unmarshalOptional[Block](result, "uncle")
}

// GetTransactionByBlockNumberAndIndex returns the transaction at the given index of a block,
// or nil if it does not exist. Equivalent to web3.eth.getTransactionFromBlock.
func (e *Eth) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNumber BlockParameter, index uint64) (*Transaction, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	result, err := e.client.Call(ctx, EthGetTransactionByBlockNumberAndIndex.String(), []interface{}{blockNumber.String(), ToHex(index)})
	if err != nil {
		return nil, err
	}
	return unmarshalOptional[Transaction](result, "transaction")
}

// GetTransactionByBlockHashAndIndex returns the transaction at the given index of the block
// with the given hash, or nil if it does not exist.
func (e *Eth) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash string, index uint64) (*Transaction, error) {
	result, err := e.client.Call(ctx, EthGetTransactionByBlockHashAndIndex.String(), []interface{}{blockHash, ToHex(index)})
	if err != nil {
		return nil, err
	}
	return unmarshalOptional[Transaction](result, "transaction")
}

// GetBlockReceipts returns all transaction receipts of a block. blockNumber may be a tag,
// a number (BlockNumber) or a block hash (BlockParameter("0x...")). Returns nil if the block is unknown.
func (e *Eth) GetBlockReceipts(ctx context.Context, blockNumber BlockParameter) ([]*TransactionReceipt, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	result, err := e.client.Call(ctx, EthGetBlockReceipts.String(), []interface{}{blockNumber.String()})
	if err != nil {
		return nil, err
	}
	var receipts []*TransactionReceipt
	if err := json.Unmarshal(result, &receipts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block receipts: %w", err)
	}
	return receipts, nil
}
