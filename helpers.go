package web3

import (
	"context"
	"fmt"
	"math/big"

	blockchainhelper "github.com/donghquinn/go-blockchain-helper/pkg/web3"
)

// Gas price helpers using go-blockchain-helper
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

// Enhanced gas estimation using go-blockchain-helper
func EstimateGasWithBuffer(ctx context.Context, client *Client, tx map[string]interface{}, buffer float64) (uint64, error) {
	baseEstimate, err := client.Eth().EstimateGas(ctx, tx)
	if err != nil {
		return 0, err
	}
	
	bufferAmount := float64(baseEstimate) * buffer
	return baseEstimate + uint64(bufferAmount), nil
}

// Unit conversion helpers using go-blockchain-helper
func EtherToWei(ether string) (*big.Int, error) {
	// Parse the ether string as float64 first
	return blockchainhelper.ParseEther(ether)
}

func WeiToEther(wei *big.Int) (string, error) {
	// Use FormatEther with default precision
	result := blockchainhelper.FormatEther(wei, 18)
	return result, nil
}

func GweiToWei(gwei string) (*big.Int, error) {
	return blockchainhelper.ParseUnits(gwei, 9) // Gwei has 9 decimals
}

func WeiToGwei(wei *big.Int) (string, error) {
	result := blockchainhelper.FormatUnits(wei, 9)
	return result, nil
}

// Enhanced unit conversion with go-blockchain-helper
func ParseEther(ether string) (*big.Int, error) {
	return blockchainhelper.ParseEther(ether)
}

func FormatEther(wei *big.Int, decimals int) string {
	return blockchainhelper.FormatEther(wei, decimals)
}

func ParseUnits(value string, decimals int) (*big.Int, error) {
	return blockchainhelper.ParseUnits(value, decimals)
}

func FormatUnits(value *big.Int, decimals int) string {
	return blockchainhelper.FormatUnits(value, decimals)
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

// Transaction helpers using go-blockchain-helper
func NewSimpleTransfer(to string, amountEth string, chainID ChainID) *TransactionParams {
	value, _ := EtherToWei(amountEth)
	return NewTransactionParams().
		SetTo(to).
		SetValue(value).
		SetGas(GasLimitTransfer.Uint64()).
		SetChainID(chainID)
}

// Enhanced transaction creation using go-blockchain-helper
func CreateTransactionWithEstimate(to string, value *big.Int, data []byte, chainID ChainID) (*TransactionParams, error) {
	// Create transaction with estimated values
	estimatedGas, err := blockchainhelper.EstimateGas("", to, "", value)
	if err != nil {
		return nil, err
	}
	suggestedGasPrice := blockchainhelper.SuggestGasPrice()
	
	return NewTransactionParams().
		SetTo(to).
		SetValue(value).
		SetData(data).
		SetGas(estimatedGas).
		SetGasPrice(suggestedGasPrice).
		SetChainID(chainID), nil
}

// ERC20 token helpers using go-blockchain-helper
func NewERC20Token(contractAddress, name, symbol string, decimals uint8) *blockchainhelper.ERC20Token {
	return blockchainhelper.NewERC20Token(contractAddress, name, symbol, decimals)
}

func EncodeERC20Transfer(token *blockchainhelper.ERC20Token, to string, amount *big.Int) ([]byte, error) {
	return token.EncodeTransfer(to, amount)
}

func EncodeERC20TransferFrom(token *blockchainhelper.ERC20Token, from, to string, amount *big.Int) ([]byte, error) {
	return token.EncodeTransferFrom(from, to, amount)
}

func EncodeERC20Approve(token *blockchainhelper.ERC20Token, spender string, amount *big.Int) ([]byte, error) {
	return token.EncodeApprove(spender, amount)
}

func NewTokenTransfer(tokenContract, to string, amount *big.Int, chainID ChainID) (*TransactionParams, error) {
	// Create a basic ERC20 token for encoding
	token := blockchainhelper.NewERC20Token(tokenContract, "Token", "TKN", 18)
	data, err := EncodeERC20Transfer(token, to, amount)
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
	// Create a basic ERC20 token for encoding
	token := blockchainhelper.NewERC20Token(tokenContract, "Token", "TKN", 18)
	data, err := EncodeERC20Approve(token, spender, amount)
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

// Enhanced contract interaction using go-blockchain-helper
func GetTokenBalance(ctx context.Context, client *Client, tokenContract, address string) (*big.Int, error) {
	token := blockchainhelper.NewERC20Token(tokenContract, "Token", "TKN", 18)
	data, err := token.EncodeBalanceOf(address)
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
	token := blockchainhelper.NewERC20Token(tokenContract, "Token", "TKN", 18)
	data, err := token.EncodeAllowance(owner, spender)
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

// Address helpers
func IsZeroAddress(address string) bool {
	return address == ZeroAddress.String() || address == "0x0"
}

func IsBurnAddress(address string) bool {
	return address == BurnAddress.String()
}

// Enhanced ABI encoding using go-blockchain-helper
func EncodeFunctionCallAdvanced(signature string, abiParams []blockchainhelper.ABIParam, params []interface{}) ([]byte, error) {
	return blockchainhelper.EncodeFunctionCall(signature, abiParams, params)
}

func DecodeFunctionResult(signatures []string, data []byte) ([]interface{}, error) {
	return blockchainhelper.DecodeFunctionResult(signatures, data)
}

// Transaction status helpers
func IsTransactionSuccess(receipt *TransactionReceipt) bool {
	return TxStatus(receipt.Status).IsSuccess()
}

func IsTransactionFailure(receipt *TransactionReceipt) bool {
	return TxStatus(receipt.Status).IsFailure()
}

// Enhanced transaction fee calculation using go-blockchain-helper
func CalculateTransactionFee(gasLimit uint64, gasPrice *big.Int) *big.Int {
	// Calculate fee manually: gasLimit * gasPrice
	gasLimitBig := new(big.Int).SetUint64(gasLimit)
	return new(big.Int).Mul(gasLimitBig, gasPrice)
}

// ERC721 helpers using go-blockchain-helper
func NewERC721Token(contractAddress, name, symbol string) *blockchainhelper.ERC721Token {
	return blockchainhelper.NewERC721Token(contractAddress, name, symbol)
}

func EncodeERC721Transfer(token *blockchainhelper.ERC721Token, from, to string, tokenId *big.Int) ([]byte, error) {
	return token.EncodeTransferFrom(from, to, tokenId)
}

func EncodeERC721Approve(token *blockchainhelper.ERC721Token, approved string, tokenId *big.Int) ([]byte, error) {
	return token.EncodeApprove(approved, tokenId)
}

func EncodeERC721SetApprovalForAll(token *blockchainhelper.ERC721Token, operator string, approved bool) ([]byte, error) {
	return token.EncodeSetApprovalForAll(operator, approved)
}

// Event processing using go-blockchain-helper
func CreateEventMonitor() *blockchainhelper.EventMonitor {
	return blockchainhelper.NewEventMonitor()
}

func ParseTransferEvent(event blockchainhelper.Event) (*blockchainhelper.TransferEvent, error) {
	return blockchainhelper.ParseTransferEvent(event)
}

// Note: ParseApprovalEvent may not be available in current version
// func ParseApprovalEvent(event blockchainhelper.Event) (*blockchainhelper.ApprovalEvent, error) {
//     return blockchainhelper.ParseApprovalEvent(event)
// }

// Crypto utilities using go-blockchain-helper (if available)
func ValidateAddress(address string) bool {
	// Use go-blockchain-helper's validation if available
	return blockchainhelper.ValidateAddress(address)
}

func PrivateKeyToAddressHelper(privateKey string) (string, error) {
	return blockchainhelper.PrivateKeyToAddress(privateKey)
}

// Note: GeneratePrivateKey may not be available in current version
// func GeneratePrivateKeyHelper() (string, error) {
//     return blockchainhelper.GeneratePrivateKey()
// }