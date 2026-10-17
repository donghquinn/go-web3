package web3

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
)

// SyncStatus is the progress object returned by eth_syncing while the node is syncing.
type SyncStatus struct {
	StartingBlock string `json:"startingBlock"`
	CurrentBlock  string `json:"currentBlock"`
	HighestBlock  string `json:"highestBlock"`
	KnownStates   string `json:"knownStates,omitempty"`
	PulledStates  string `json:"pulledStates,omitempty"`
}

// IsSyncing returns the node's sync progress, or nil when the node is not syncing.
func (e *Eth) IsSyncing(ctx context.Context) (*SyncStatus, error) {
	result, err := e.client.Call(ctx, EthSyncing.String(), []interface{}{})
	if err != nil {
		return nil, err
	}
	var syncing bool
	if err := json.Unmarshal(result, &syncing); err == nil {
		// eth_syncing answers a bare false when fully synced.
		return nil, nil
	}
	var status SyncStatus
	if err := json.Unmarshal(result, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sync status: %w", err)
	}
	return &status, nil
}

// parseJSONBool unmarshals a JSON-encoded boolean value.
func parseJSONBool(result json.RawMessage, field string) (bool, error) {
	var val bool
	if err := json.Unmarshal(result, &val); err != nil {
		return false, fmt.Errorf("failed to unmarshal %s: %w", field, err)
	}
	return val, nil
}

// GetCoinbase returns the address that receives the node's mining/validation rewards.
func (e *Eth) GetCoinbase(ctx context.Context) (string, error) {
	result, err := e.client.Call(ctx, EthCoinbase.String(), []interface{}{})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "coinbase")
}

// IsMining reports whether the node is mining. Always false on proof-of-stake networks.
func (e *Eth) IsMining(ctx context.Context) (bool, error) {
	result, err := e.client.Call(ctx, EthMining.String(), []interface{}{})
	if err != nil {
		return false, err
	}
	return parseJSONBool(result, "mining")
}

// GetHashRate returns the number of hashes per second the node is mining with.
func (e *Eth) GetHashRate(ctx context.Context) (*big.Int, error) {
	result, err := e.client.Call(ctx, EthHashrate.String(), []interface{}{})
	if err != nil {
		return nil, err
	}
	return parseHexBigInt(result, "hashrate")
}

// GetAccounts returns the addresses managed by the node.
func (e *Eth) GetAccounts(ctx context.Context) ([]string, error) {
	result, err := e.client.Call(ctx, EthAccounts.String(), []interface{}{})
	if err != nil {
		return nil, err
	}
	var accounts []string
	if err := json.Unmarshal(result, &accounts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal accounts: %w", err)
	}
	return accounts, nil
}

// GetProtocolVersion returns the Ethereum protocol version of the node.
func (e *Eth) GetProtocolVersion(ctx context.Context) (string, error) {
	result, err := e.client.Call(ctx, EthProtocolVersion.String(), []interface{}{})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "protocol version")
}

// GetBlobBaseFee returns the current base fee per blob gas (EIP-4844).
func (e *Eth) GetBlobBaseFee(ctx context.Context) (*big.Int, error) {
	result, err := e.client.Call(ctx, EthBlobBaseFee.String(), []interface{}{})
	if err != nil {
		return nil, err
	}
	return parseHexBigInt(result, "blob base fee")
}

// IsListening reports whether the node is listening for network connections.
func (e *Eth) IsListening(ctx context.Context) (bool, error) {
	result, err := e.client.Call(ctx, NetListening.String(), []interface{}{})
	if err != nil {
		return false, err
	}
	return parseJSONBool(result, "listening")
}

// GetPeerCount returns the number of peers connected to the node.
func (e *Eth) GetPeerCount(ctx context.Context) (uint64, error) {
	result, err := e.client.Call(ctx, NetPeerCount.String(), []interface{}{})
	if err != nil {
		return 0, err
	}
	return parseHexUint64(result, "peer count")
}

// Sha3 returns the Keccak-256 hash of hex-encoded data, computed by the node.
func (e *Eth) Sha3(ctx context.Context, data string) (string, error) {
	result, err := e.client.Call(ctx, Web3Sha3.String(), []interface{}{data})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "sha3")
}
