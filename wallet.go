package web3

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
)

type Wallet struct {
	privateKey *ecdsa.PrivateKey
	address    string
	client     *Client
}

type TransferOptions struct {
	To       string
	Value    *big.Int
	GasLimit uint64
	GasPrice *big.Int
	Data     []byte
}

type SendTransactionResult struct {
	TransactionHash string
	From            string
	To              string
	Value           *big.Int
	GasUsed         uint64
	BlockNumber     uint64
	Status          bool
}

func NewWallet(privateKeyHex string, client *Client) (*Wallet, error) {
	privateKey, err := PrivateKeyFromHex(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	return &Wallet{
		privateKey: privateKey,
		address:    PrivateKeyToAddress(privateKey),
		client:     client,
	}, nil
}

func CreateWallet(client *Client) (*Wallet, error) {
	privateKey, err := GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return &Wallet{
		privateKey: privateKey,
		address:    PrivateKeyToAddress(privateKey),
		client:     client,
	}, nil
}

func (w *Wallet) GetAddress() string    { return w.address }
func (w *Wallet) GetPrivateKey() string { return PrivateKeyToHex(w.privateKey) }

func (w *Wallet) GetBalance(ctx context.Context) (*big.Int, error) {
	return w.client.Eth().GetBalance(ctx, w.address, BlockLatest)
}

func (w *Wallet) GetNonce(ctx context.Context) (uint64, error) {
	return w.client.Eth().GetTransactionCount(ctx, w.address, BlockPending)
}

// resolveGasLimit estimates and sets opts.GasLimit when it is zero.
func (w *Wallet) resolveGasLimit(ctx context.Context, opts *TransferOptions, bufferPct uint64) error {
	if opts.GasLimit > 0 {
		return nil
	}
	callObj := map[string]interface{}{
		"from":  w.address,
		"to":    opts.To,
		"value": fmt.Sprintf("0x%x", opts.Value),
		"data":  fmt.Sprintf("0x%x", opts.Data),
	}
	estimate, err := w.client.Eth().EstimateGas(ctx, callObj)
	if err != nil {
		return fmt.Errorf("failed to estimate gas: %w", err)
	}
	opts.GasLimit = estimate + (estimate * bufferPct / 100)
	return nil
}

// resolveGasPrice fetches and sets opts.GasPrice when it is nil.
func (w *Wallet) resolveGasPrice(ctx context.Context, opts *TransferOptions) error {
	if opts.GasPrice != nil {
		return nil
	}
	price, err := w.client.Eth().GetGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}
	opts.GasPrice = price
	return nil
}

func (w *Wallet) SendTransaction(ctx context.Context, opts *TransferOptions) (*SendTransactionResult, error) {
	if err := w.resolveGasLimit(ctx, opts, 10); err != nil {
		return nil, err
	}
	if err := w.resolveGasPrice(ctx, opts); err != nil {
		return nil, err
	}

	nonce, err := w.GetNonce(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	txParams := NewTransactionParams().
		SetTo(opts.To).
		SetValue(opts.Value).
		SetGas(opts.GasLimit).
		SetGasPrice(opts.GasPrice).
		SetData(opts.Data).
		SetNonce(nonce).
		SetChainID(ChainMainnet)

	signedTx, err := SignTransaction(txParams, w.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	txHash, err := w.client.Eth().SendRawTransaction(ctx, signedTx.Raw)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return &SendTransactionResult{
		TransactionHash: txHash,
		From:            w.address,
		To:              opts.To,
		Value:           opts.Value,
	}, nil
}

func (w *Wallet) SendEther(ctx context.Context, to, amountInEther string) (*SendTransactionResult, error) {
	value, err := ToWei(amountInEther, Ether)
	if err != nil {
		return nil, fmt.Errorf("invalid ether amount: %w", err)
	}
	return w.SendTransaction(ctx, &TransferOptions{To: to, Value: value})
}

func (w *Wallet) SendWei(ctx context.Context, to string, amountInWei *big.Int) (*SendTransactionResult, error) {
	return w.SendTransaction(ctx, &TransferOptions{To: to, Value: amountInWei})
}

func (w *Wallet) SendEIP1559Transaction(ctx context.Context, opts *TransferOptions, maxFeePerGas, maxPriorityFeePerGas *big.Int) (*SendTransactionResult, error) {
	if err := w.resolveGasLimit(ctx, opts, 10); err != nil {
		return nil, err
	}

	nonce, err := w.GetNonce(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	txParams := NewEIP1559TransactionParams()
	txParams.To = opts.To
	txParams.Value = opts.Value
	txParams.Gas = opts.GasLimit
	txParams.MaxFeePerGas = maxFeePerGas
	txParams.MaxPriorityFeePerGas = maxPriorityFeePerGas
	txParams.Data = opts.Data
	txParams.Nonce = nonce
	txParams.ChainID = ChainMainnet.BigInt()

	signedTx, err := SignEIP1559Transaction(txParams, w.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	txHash, err := w.client.Eth().SendRawTransaction(ctx, signedTx.Raw)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return &SendTransactionResult{
		TransactionHash: txHash,
		From:            w.address,
		To:              opts.To,
		Value:           opts.Value,
	}, nil
}

func (w *Wallet) CallContract(ctx context.Context, contractAddress string, methodData []byte) (string, error) {
	callObj := map[string]interface{}{
		"from": w.address,
		"to":   contractAddress,
		"data": fmt.Sprintf("0x%x", methodData),
	}
	return w.client.Eth().Call(ctx, callObj, BlockLatest)
}

func (w *Wallet) SendContractTransaction(ctx context.Context, contractAddress string, methodData []byte, value *big.Int) (*SendTransactionResult, error) {
	return w.SendTransaction(ctx, &TransferOptions{
		To:    contractAddress,
		Value: value,
		Data:  methodData,
	})
}

func (w *Wallet) DeployContract(ctx context.Context, bytecode, constructorData []byte, gasLimit uint64, gasPrice *big.Int) (*SendTransactionResult, error) {
	deployData := append(bytecode, constructorData...)
	opts := &TransferOptions{
		Value:    big.NewInt(0),
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		Data:     deployData,
	}

	if opts.GasLimit == 0 {
		callObj := map[string]interface{}{
			"from": w.address,
			"data": fmt.Sprintf("0x%x", deployData),
		}
		estimate, err := w.client.Eth().EstimateGas(ctx, callObj)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas: %w", err)
		}
		opts.GasLimit = estimate + (estimate * 20 / 100)
	}
	if err := w.resolveGasPrice(ctx, opts); err != nil {
		return nil, err
	}

	return w.SendTransaction(ctx, opts)
}

func (w *Wallet) WaitForTransaction(ctx context.Context, txHash string) (*TransactionReceipt, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			receipt, err := w.client.Eth().GetTransactionReceipt(ctx, txHash)
			if err == nil {
				return receipt, nil
			}
		}
	}
}
