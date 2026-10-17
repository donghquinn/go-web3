package web3

import (
	"context"
	"encoding/json"
	"fmt"
)

// NewFilter installs a log filter on the node and returns its ID.
// Poll it with GetFilterLogChanges and remove it with UninstallFilter.
func (e *Eth) NewFilter(ctx context.Context, filter LogFilter) (string, error) {
	result, err := e.client.Call(ctx, EthNewFilter.String(), []interface{}{filter})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "filter ID")
}

// NewBlockFilter installs a filter that reports new block hashes. Poll it with GetFilterHashChanges.
func (e *Eth) NewBlockFilter(ctx context.Context) (string, error) {
	result, err := e.client.Call(ctx, EthNewBlockFilter.String(), []interface{}{})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "filter ID")
}

// NewPendingTransactionFilter installs a filter that reports new pending transaction hashes.
// Poll it with GetFilterHashChanges.
func (e *Eth) NewPendingTransactionFilter(ctx context.Context) (string, error) {
	result, err := e.client.Call(ctx, EthNewPendingTransactionFilter.String(), []interface{}{})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "filter ID")
}

// GetFilterLogChanges returns the logs matched by a log filter since the last poll.
func (e *Eth) GetFilterLogChanges(ctx context.Context, filterID string) ([]*Log, error) {
	result, err := e.client.Call(ctx, EthGetFilterChanges.String(), []interface{}{filterID})
	if err != nil {
		return nil, err
	}
	var logs []*Log
	if err := json.Unmarshal(result, &logs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter changes: %w", err)
	}
	return logs, nil
}

// GetFilterHashChanges returns the block or transaction hashes reported by a block or
// pending-transaction filter since the last poll.
func (e *Eth) GetFilterHashChanges(ctx context.Context, filterID string) ([]string, error) {
	result, err := e.client.Call(ctx, EthGetFilterChanges.String(), []interface{}{filterID})
	if err != nil {
		return nil, err
	}
	var hashes []string
	if err := json.Unmarshal(result, &hashes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter changes: %w", err)
	}
	return hashes, nil
}

// GetFilterLogs returns all logs matching a log filter.
func (e *Eth) GetFilterLogs(ctx context.Context, filterID string) ([]*Log, error) {
	result, err := e.client.Call(ctx, EthGetFilterLogs.String(), []interface{}{filterID})
	if err != nil {
		return nil, err
	}
	var logs []*Log
	if err := json.Unmarshal(result, &logs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter logs: %w", err)
	}
	return logs, nil
}

// UninstallFilter removes a filter from the node. Returns false if the filter did not exist.
func (e *Eth) UninstallFilter(ctx context.Context, filterID string) (bool, error) {
	result, err := e.client.Call(ctx, EthUninstallFilter.String(), []interface{}{filterID})
	if err != nil {
		return false, err
	}
	return parseJSONBool(result, "uninstall result")
}
