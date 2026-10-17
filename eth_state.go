package web3

import (
	"context"
	"encoding/json"
	"fmt"
)

// StorageProof is the Merkle proof of a single storage slot (EIP-1186).
type StorageProof struct {
	Key   string   `json:"key"`
	Value string   `json:"value"`
	Proof []string `json:"proof"`
}

// AccountProof is the result of eth_getProof (EIP-1186).
type AccountProof struct {
	Address      string          `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      string          `json:"balance"`
	CodeHash     string          `json:"codeHash"`
	Nonce        string          `json:"nonce"`
	StorageHash  string          `json:"storageHash"`
	StorageProof []*StorageProof `json:"storageProof"`
}

// GetProof returns the account and storage values of an address, including Merkle proofs.
func (e *Eth) GetProof(ctx context.Context, address string, storageKeys []string, blockNumber BlockParameter) (*AccountProof, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	if storageKeys == nil {
		storageKeys = []string{}
	}
	result, err := e.client.Call(ctx, EthGetProof.String(), []interface{}{address, storageKeys, blockNumber.String()})
	if err != nil {
		return nil, err
	}
	var proof AccountProof
	if err := json.Unmarshal(result, &proof); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proof: %w", err)
	}
	return &proof, nil
}

// AccessListResult is the result of eth_createAccessList (EIP-2930).
type AccessListResult struct {
	AccessList []AccessTuple `json:"accessList"`
	GasUsed    string        `json:"gasUsed"`
	Error      string        `json:"error,omitempty"` // set when the simulated call reverts
}

// CreateAccessList simulates a transaction and returns the access list it would touch.
func (e *Eth) CreateAccessList(ctx context.Context, callObj map[string]interface{}, blockNumber BlockParameter) (*AccessListResult, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}
	result, err := e.client.Call(ctx, EthCreateAccessList.String(), []interface{}{callObj, blockNumber.String()})
	if err != nil {
		return nil, err
	}
	var accessList AccessListResult
	if err := json.Unmarshal(result, &accessList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access list: %w", err)
	}
	return &accessList, nil
}

// Sign signs hex-encoded data with an account unlocked on the node (eth_sign).
func (e *Eth) Sign(ctx context.Context, address, data string) (string, error) {
	result, err := e.client.Call(ctx, EthSign.String(), []interface{}{address, data})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "signature")
}

// SignTransactionResult is the result of eth_signTransaction.
type SignTransactionResult struct {
	Raw string       `json:"raw"`
	Tx  *Transaction `json:"tx"`
}

// SignTransaction signs a transaction with an account unlocked on the node, without sending it.
func (e *Eth) SignTransaction(ctx context.Context, tx map[string]interface{}) (*SignTransactionResult, error) {
	result, err := e.client.Call(ctx, EthSignTransaction.String(), []interface{}{tx})
	if err != nil {
		return nil, err
	}
	var signed SignTransactionResult
	if err := json.Unmarshal(result, &signed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal signed transaction: %w", err)
	}
	return &signed, nil
}

// SendTransaction signs and sends a transaction with an account unlocked on the node.
// To send from a local private key, use Wallet or SendRawTransaction instead.
func (e *Eth) SendTransaction(ctx context.Context, tx map[string]interface{}) (string, error) {
	result, err := e.client.Call(ctx, EthSendTransaction.String(), []interface{}{tx})
	if err != nil {
		return "", err
	}
	return parseJSONString(result, "transaction hash")
}
