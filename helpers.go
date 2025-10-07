package web3

import (
	"context"
	"fmt"
	"math/big"
)

// Gas price helpers
func GetOptimalGasPrice(ctx context.Context, client *Client, level GasPriceLevel) (*big.Int, error) {
	basePrice, err := client.Eth().GetGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	
	multiplier := level.Multiplier()
	factor := new(big.Float).SetFloat64(multiplier)
	
	result := new(big.Float).Mul(new(big.Float).SetInt(basePrice), factor)
	optimal, _ := result.Int(nil)
	
	return optimal, nil
}

// Chain helpers
func GetNetworkConfig(chainID ChainID) (NetworkConfig, error) {
	config, exists := Networks[chainID]
	if !exists {
		return NetworkConfig{}, fmt.Errorf("network configuration not found for chain ID %d", chainID)
	}
	return config, nil
}

func IsTestnet(chainID ChainID) bool {
	testnets := []ChainID{
		ChainGoerli, ChainSepolia, ChainOptimismGoerli, 
		ChainArbitrumGoerli, ChainPolygonMumbai, 
		ChainAvalancheFuji, ChainBSCTestnet, ChainFantomTestnet,
	}
	
	for _, testnet := range testnets {
		if chainID == testnet {
			return true
		}
	}
	return false
}

func IsMainnet(chainID ChainID) bool {
	return chainID == ChainMainnet
}

// Unit conversion helpers
func EtherToWei(ether string) (*big.Int, error) {
	return ToWei(ether, Ether)
}

func WeiToEther(wei *big.Int) (string, error) {
	return FromWei(wei, Ether)
}

func GweiToWei(gwei string) (*big.Int, error) {
	return ToWei(gwei, Gwei)
}

func WeiToGwei(wei *big.Int) (string, error) {
	return FromWei(wei, Gwei)
}

// Transaction helpers
func NewSimpleTransfer(to string, amountEth string, chainID ChainID) *TransactionParams {
	value, _ := EtherToWei(amountEth)
	return NewTransactionParams().
		SetTo(to).
		SetValue(value).
		SetGas(GasLimitTransfer.Uint64()).
		SetChainID(chainID)
}

func NewTokenTransfer(tokenContract, to string, amount *big.Int, chainID ChainID) (*TransactionParams, error) {
	data, err := EncodeABI(FuncTransfer.String(), to, amount)
	if err != nil {
		return nil, err
	}
	
	return NewTransactionParams().
		SetTo(tokenContract).
		SetValue(big.NewInt(0)).
		SetData(data).
		SetGas(GasLimitTokenTransfer.Uint64()).
		SetChainID(chainID), nil
}

func NewTokenApproval(tokenContract, spender string, amount *big.Int, chainID ChainID) (*TransactionParams, error) {
	data, err := EncodeABI(FuncApprove.String(), spender, amount)
	if err != nil {
		return nil, err
	}
	
	return NewTransactionParams().
		SetTo(tokenContract).
		SetValue(big.NewInt(0)).
		SetData(data).
		SetGas(GasLimitTokenApproval.Uint64()).
		SetChainID(chainID), nil
}

// Address helpers
func IsZeroAddress(address string) bool {
	return address == ZeroAddress.String() || address == "0x0"
}

func IsBurnAddress(address string) bool {
	return address == BurnAddress.String()
}

// Common contract calls
func GetTokenBalance(ctx context.Context, client *Client, tokenContract, address string) (*big.Int, error) {
	data, err := EncodeABI(FuncBalanceOf.String(), address)
	if err != nil {
		return nil, err
	}
	
	callObj := map[string]interface{}{
		"to":   tokenContract,
		"data": fmt.Sprintf("0x%x", data),
	}
	
	result, err := client.Eth().Call(ctx, callObj, BlockLatest)
	if err != nil {
		return nil, err
	}
	
	return FromHex(result)
}

func GetTokenAllowance(ctx context.Context, client *Client, tokenContract, owner, spender string) (*big.Int, error) {
	data, err := EncodeABI(FuncAllowance.String(), owner, spender)
	if err != nil {
		return nil, err
	}
	
	callObj := map[string]interface{}{
		"to":   tokenContract,
		"data": fmt.Sprintf("0x%x", data),
	}
	
	result, err := client.Eth().Call(ctx, callObj, BlockLatest)
	if err != nil {
		return nil, err
	}
	
	return FromHex(result)
}

// Transaction status helpers
func IsTransactionSuccess(receipt *TransactionReceipt) bool {
	return TxStatus(receipt.Status).IsSuccess()
}

func IsTransactionFailure(receipt *TransactionReceipt) bool {
	return TxStatus(receipt.Status).IsFailure()
}

// Gas estimation with buffers
func EstimateGasWithBuffer(ctx context.Context, client *Client, tx map[string]interface{}, buffer float64) (uint64, error) {
	baseEstimate, err := client.Eth().EstimateGas(ctx, tx)
	if err != nil {
		return 0, err
	}
	
	bufferAmount := float64(baseEstimate) * buffer
	return baseEstimate + uint64(bufferAmount), nil
}