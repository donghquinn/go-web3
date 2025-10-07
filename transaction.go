package web3

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	blockchainhelper "github.com/donghquinn/go-blockchain-helper/pkg/web3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type TransactionParams struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	Value    *big.Int `json:"value"`
	Gas      uint64   `json:"gas"`
	GasPrice *big.Int `json:"gasPrice"`
	Data     []byte   `json:"data"`
	Nonce    uint64   `json:"nonce"`
	ChainID  *big.Int `json:"chainId"`
}

type EIP1559TransactionParams struct {
	From                 string   `json:"from"`
	To                   string   `json:"to"`
	Value                *big.Int `json:"value"`
	Gas                  uint64   `json:"gas"`
	MaxFeePerGas         *big.Int `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas"`
	Data                 []byte   `json:"data"`
	Nonce                uint64   `json:"nonce"`
	ChainID              *big.Int `json:"chainId"`
}

type SignedTransaction struct {
	Hash string `json:"hash"`
	Raw  string `json:"raw"`
}

func NewTransactionParams() *TransactionParams {
	return &TransactionParams{
		Value:   big.NewInt(0),
		Data:    []byte{},
		ChainID: ChainMainnet.BigInt(),
	}
}

func NewEIP1559TransactionParams() *EIP1559TransactionParams {
	return &EIP1559TransactionParams{
		Value:   big.NewInt(0),
		Data:    []byte{},
		ChainID: ChainMainnet.BigInt(),
	}
}

func (tp *TransactionParams) SetTo(address string) *TransactionParams {
	tp.To = address
	return tp
}

func (tp *TransactionParams) SetValue(value *big.Int) *TransactionParams {
	tp.Value = value
	return tp
}

func (tp *TransactionParams) SetValueInWei(wei string) *TransactionParams {
	value, _ := new(big.Int).SetString(wei, 10)
	tp.Value = value
	return tp
}

func (tp *TransactionParams) SetValueInEther(eth string) *TransactionParams {
	value, _ := ToWei(eth, Ether)
	tp.Value = value
	return tp
}

func (tp *TransactionParams) SetGas(gas uint64) *TransactionParams {
	tp.Gas = gas
	return tp
}

func (tp *TransactionParams) SetGasPrice(gasPrice *big.Int) *TransactionParams {
	tp.GasPrice = gasPrice
	return tp
}

func (tp *TransactionParams) SetGasPriceInGwei(gwei string) *TransactionParams {
	gasPrice, _ := ToWei(gwei, Gwei)
	tp.GasPrice = gasPrice
	return tp
}

func (tp *TransactionParams) SetData(data []byte) *TransactionParams {
	tp.Data = data
	return tp
}

func (tp *TransactionParams) SetDataFromHex(hexData string) *TransactionParams {
	if len(hexData) >= 2 && hexData[:2] == "0x" {
		hexData = hexData[2:]
	}
	data, _ := hex.DecodeString(hexData)
	tp.Data = data
	return tp
}

func (tp *TransactionParams) SetNonce(nonce uint64) *TransactionParams {
	tp.Nonce = nonce
	return tp
}

func (tp *TransactionParams) SetChainID(chainID ChainID) *TransactionParams {
	tp.ChainID = chainID.BigInt()
	return tp
}

func PrivateKeyFromHex(hexKey string) (*ecdsa.PrivateKey, error) {
	if len(hexKey) >= 2 && hexKey[:2] == "0x" {
		hexKey = hexKey[2:]
	}
	
	privateKeyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}
	
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	
	return privateKey, nil
}

func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

func PrivateKeyToHex(privateKey *ecdsa.PrivateKey) string {
	return fmt.Sprintf("0x%x", crypto.FromECDSA(privateKey))
}

func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) string {
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey)
	return address.Hex()
}

func SignTransaction(tx *TransactionParams, privateKey *ecdsa.PrivateKey) (*SignedTransaction, error) {
	if tx.To == "" {
		return nil, fmt.Errorf("transaction recipient (to) is required")
	}
	if tx.GasPrice == nil {
		return nil, fmt.Errorf("gas price is required")
	}
	if tx.Gas == 0 {
		return nil, fmt.Errorf("gas limit is required")
	}

	var toAddr *common.Address
	if tx.To != "" {
		addr := common.HexToAddress(tx.To)
		toAddr = &addr
	}

	ethTx := types.NewTx(&types.LegacyTx{
		Nonce:    tx.Nonce,
		To:       toAddr,
		Value:    tx.Value,
		Gas:      tx.Gas,
		GasPrice: tx.GasPrice,
		Data:     tx.Data,
	})

	signer := types.NewEIP155Signer(tx.ChainID)
	signedTx, err := types.SignTx(ethTx, signer, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	rawTxBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to encode transaction: %w", err)
	}

	return &SignedTransaction{
		Hash: signedTx.Hash().Hex(),
		Raw:  fmt.Sprintf("0x%x", rawTxBytes),
	}, nil
}

func SignEIP1559Transaction(tx *EIP1559TransactionParams, privateKey *ecdsa.PrivateKey) (*SignedTransaction, error) {
	if tx.To == "" {
		return nil, fmt.Errorf("transaction recipient (to) is required")
	}
	if tx.MaxFeePerGas == nil {
		return nil, fmt.Errorf("maxFeePerGas is required")
	}
	if tx.MaxPriorityFeePerGas == nil {
		return nil, fmt.Errorf("maxPriorityFeePerGas is required")
	}
	if tx.Gas == 0 {
		return nil, fmt.Errorf("gas limit is required")
	}

	var toAddr *common.Address
	if tx.To != "" {
		addr := common.HexToAddress(tx.To)
		toAddr = &addr
	}

	ethTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   tx.ChainID,
		Nonce:     tx.Nonce,
		To:        toAddr,
		Value:     tx.Value,
		Gas:       tx.Gas,
		GasTipCap: tx.MaxPriorityFeePerGas,
		GasFeeCap: tx.MaxFeePerGas,
		Data:      tx.Data,
	})

	signer := types.NewLondonSigner(tx.ChainID)
	signedTx, err := types.SignTx(ethTx, signer, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	rawTxBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to encode transaction: %w", err)
	}

	return &SignedTransaction{
		Hash: signedTx.Hash().Hex(),
		Raw:  fmt.Sprintf("0x%x", rawTxBytes),
	}, nil
}

func CreateContractDeployment(bytecode []byte, constructorData []byte, privateKey *ecdsa.PrivateKey, params *TransactionParams) (*SignedTransaction, error) {
	params.To = ""
	
	if constructorData != nil {
		params.Data = append(bytecode, constructorData...)
	} else {
		params.Data = bytecode
	}

	return SignTransaction(params, privateKey)
}

func CreateContractCall(contractAddress string, methodData []byte, privateKey *ecdsa.PrivateKey, params *TransactionParams) (*SignedTransaction, error) {
	params.To = contractAddress
	params.Data = methodData

	return SignTransaction(params, privateKey)
}

func RecoverSigner(rawTxHex string) (string, error) {
	if len(rawTxHex) >= 2 && rawTxHex[:2] == "0x" {
		rawTxHex = rawTxHex[2:]
	}

	rawTxBytes, err := hex.DecodeString(rawTxHex)
	if err != nil {
		return "", fmt.Errorf("invalid hex string: %w", err)
	}

	var tx types.Transaction
	err = rlp.DecodeBytes(rawTxBytes, &tx)
	if err != nil {
		return "", fmt.Errorf("failed to decode transaction: %w", err)
	}

	var signer types.Signer
	if tx.ChainId().Cmp(big.NewInt(0)) == 0 {
		signer = types.HomesteadSigner{}
	} else {
		signer = types.NewEIP155Signer(tx.ChainId())
	}

	sender, err := signer.Sender(&tx)
	if err != nil {
		return "", fmt.Errorf("failed to recover sender: %w", err)
	}

	return sender.Hex(), nil
}

func EncodeABI(methodSignature string, params ...interface{}) ([]byte, error) {
	// Convert params to slice for go-blockchain-helper
	paramSlice := make([]interface{}, len(params))
	copy(paramSlice, params)
	
	// Create basic ABI params - this is a simplified approach
	// In a real implementation, you would parse the method signature to determine types
	abiParams := make([]blockchainhelper.ABIParam, len(params))
	for i, param := range params {
		switch param.(type) {
		case string:
			if IsAddress(param.(string)) {
				abiParams[i] = blockchainhelper.ABIParam{Type: "address"}
			} else {
				abiParams[i] = blockchainhelper.ABIParam{Type: "string"}
			}
		case *big.Int:
			abiParams[i] = blockchainhelper.ABIParam{Type: "uint256"}
		case uint64:
			abiParams[i] = blockchainhelper.ABIParam{Type: "uint64"}
		case []byte:
			abiParams[i] = blockchainhelper.ABIParam{Type: "bytes"}
		case bool:
			abiParams[i] = blockchainhelper.ABIParam{Type: "bool"}
		default:
			return nil, fmt.Errorf("unsupported parameter type: %T", param)
		}
	}
	
	// Use go-blockchain-helper for ABI encoding
	return blockchainhelper.EncodeFunctionCall(methodSignature, abiParams, paramSlice)
}

func RandomNonce() uint64 {
	nonce := make([]byte, 8)
	rand.Read(nonce)
	return new(big.Int).SetBytes(nonce).Uint64()
}