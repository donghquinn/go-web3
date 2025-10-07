# Go-Web3

A comprehensive Go library for interacting with Ethereum blockchain through JSON-RPC APIs, designed with web3.js-like interface for familiar usage patterns.

## 🚀 Features

- **Complete Ethereum JSON-RPC Implementation**: All major eth_* methods supported
- **Context-Aware Operations**: Proper context handling for timeouts and cancellation
- **Type-Safe Structures**: Strongly-typed transaction, block, and receipt objects
- **Web3.js-Like API**: Familiar method names and usage patterns for JavaScript developers
- **Built-in Utilities**: Wei/Ether conversion, address validation, hex operations
- **Robust Error Handling**: Detailed RPC error information with proper Go error wrapping
- **Production Ready**: Thread-safe client with connection pooling

## 📦 Installation

```bash
go get github.com/donghquinn/go-web3
```

## 🏃‍♂️ Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/donghquinn/go-web3"
)

func main() {
    // Create a new client with your Ethereum RPC endpoint
    client := web3.NewClient("https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY")
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Get latest block number
    blockNumber, err := client.Eth().GetBlockNumber(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Latest block: %d\n", blockNumber)
    
    // Check account balance
    address := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
    balance, err := client.Eth().GetBalance(ctx, address, "latest")
    if err != nil {
        log.Fatal(err)
    }
    
    // Convert Wei to Ether for display
    balanceEth, _ := web3.FromWei(balance, "ether")
    fmt.Printf("Balance: %s ETH\n", balanceEth)
}
```

## 📖 Detailed Usage Guide

### Client Configuration

#### Basic Client
```go
client := web3.NewClient("https://mainnet.infura.io/v3/YOUR_PROJECT_ID")
```

#### Client with Custom HTTP Client
```go
import (
    "net/http"
    "time"
)

httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
    },
}

client := web3.NewClient("https://mainnet.infura.io/v3/YOUR_PROJECT_ID")
// Note: Custom HTTP client configuration would require extending the library
```

### Context Usage

Always use context for proper timeout and cancellation handling:

```go
// Context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Context with deadline
deadline := time.Now().Add(1 * time.Minute)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// Cancellable context
ctx, cancel := context.WithCancel(context.Background())
// Call cancel() when needed
```

## 📚 Complete API Reference

### Client Methods

#### Creating a Client
```go
client := web3.NewClient("YOUR_ETHEREUM_RPC_URL")

// Popular RPC endpoints:
// Mainnet: "https://mainnet.infura.io/v3/YOUR_PROJECT_ID"
// Goerli: "https://goerli.infura.io/v3/YOUR_PROJECT_ID"
// Alchemy: "https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY"
// Local: "http://localhost:8545"
```

### Ethereum Methods (client.Eth())

#### 💰 Account & Balance Operations

##### Get Balance
```go
// Get balance in Wei
balance, err := client.Eth().GetBalance(ctx, address, "latest")
if err != nil {
    log.Fatal(err)
}

// Convert to Ether for display
balanceEth, _ := web3.FromWei(balance, "ether")
fmt.Printf("Balance: %s ETH\n", balanceEth)

// Block parameters: "latest", "earliest", "pending", or hex block number
balance, err := client.Eth().GetBalance(ctx, address, "0x1b4") // specific block
```

##### Get Transaction Count (Nonce)
```go
// Get nonce for transaction
nonce, err := client.Eth().GetTransactionCount(ctx, address, "latest")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Next nonce: %d\n", nonce)

// For pending transactions (useful for rapid transaction sending)
pendingNonce, err := client.Eth().GetTransactionCount(ctx, address, "pending")
```

#### 🔗 Block Operations

##### Get Current Block Number
```go
blockNumber, err := client.Eth().GetBlockNumber(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Latest block: %d\n", blockNumber)
```

##### Get Block by Number
```go
// Get block without full transaction details
block, err := client.Eth().GetBlockByNumber(ctx, "latest", false)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Block Hash: %s\n", block.Hash)
fmt.Printf("Miner: %s\n", block.Miner)
fmt.Printf("Gas Used: %s\n", block.GasUsed)
fmt.Printf("Timestamp: %s\n", block.Timestamp)
fmt.Printf("Transaction count: %d\n", len(block.Transactions))

// Get block with full transaction details
blockWithTxs, err := client.Eth().GetBlockByNumber(ctx, "latest", true)
```

##### Get Block by Hash
```go
blockHash := "0x1234567890abcdef..."
block, err := client.Eth().GetBlockByHash(ctx, blockHash, false)
if err != nil {
    log.Fatal(err)
}
```

#### 💸 Transaction Operations

##### Get Transaction by Hash
```go
txHash := "0xabcdef1234567890..."
tx, err := client.Eth().GetTransactionByHash(ctx, txHash)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("From: %s\n", tx.From)
fmt.Printf("To: %s\n", tx.To)
fmt.Printf("Value: %s\n", tx.Value)
fmt.Printf("Gas: %s\n", tx.Gas)
fmt.Printf("Gas Price: %s\n", tx.GasPrice)
```

##### Get Transaction Receipt
```go
txHash := "0xabcdef1234567890..."
receipt, err := client.Eth().GetTransactionReceipt(ctx, txHash)
if err != nil {
    log.Fatal(err)
}

// Check if transaction was successful
if receipt.Status == "0x1" {
    fmt.Println("Transaction successful")
} else {
    fmt.Println("Transaction failed")
}

fmt.Printf("Gas Used: %s\n", receipt.GasUsed)
fmt.Printf("Block Number: %s\n", receipt.BlockNumber)
```

##### Send Raw Transaction
```go
// Send a pre-signed transaction
signedTxHex := "0xf86c808504a817c800825208940xd8da6bf26964af9d7eed9e03e53415d37aa96045880de0b6b3a764000080820a95a0..."
txHash, err := client.Eth().SendRawTransaction(ctx, signedTxHex)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Transaction sent: %s\n", txHash)
```

#### ⛽ Gas Operations

##### Get Gas Price
```go
gasPrice, err := client.Eth().GetGasPrice(ctx)
if err != nil {
    log.Fatal(err)
}

// Convert to Gwei for display
gasPriceGwei, _ := web3.FromWei(gasPrice, "gwei")
fmt.Printf("Current gas price: %s Gwei\n", gasPriceGwei)
```

##### Estimate Gas
```go
// Create transaction object
txObj := map[string]interface{}{
    "from":  "0x742d35Cc6084C0532C9d2b908B8C0c9ff3e3ba0A",
    "to":    "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
    "value": "0xde0b6b3a7640000", // 1 ETH in wei
}

gasEstimate, err := client.Eth().EstimateGas(ctx, txObj)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Estimated gas: %d\n", gasEstimate)
```

#### 📞 Contract Interaction

##### Call Contract Method
```go
// Call a contract method (read-only)
callObj := map[string]interface{}{
    "to":   "0xA0b86a33E6417c48cd7a94Ca95e70aD2c51e74f7", // contract address
    "data": "0x70a08231000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045", // balanceOf call
}

result, err := client.Eth().Call(ctx, callObj, "latest")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Contract call result: %s\n", result)
```

## 💰 Transaction Management

### Creating and Managing Wallets

#### Generate New Wallet
```go
// Create a new random wallet
wallet, err := web3.CreateWallet(client)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Address: %s\n", wallet.GetAddress())
fmt.Printf("Private Key: %s\n", wallet.GetPrivateKey()) // Keep this secure!
```

#### Load Existing Wallet
```go
// Load wallet from private key
privateKey := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
wallet, err := web3.NewWallet(privateKey, client)
if err != nil {
    log.Fatal(err)
}

// Check wallet balance
balance, err := wallet.GetBalance(ctx)
if err != nil {
    log.Fatal(err)
}

balanceEth, _ := web3.FromWei(balance, "ether")
fmt.Printf("Balance: %s ETH\n", balanceEth)
```

### Building and Signing Transactions

#### Legacy Transaction (Pre-EIP-1559)
```go
// Create transaction parameters
txParams := web3.NewTransactionParams().
    SetTo("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045").
    SetValueInEther("0.1").           // Send 0.1 ETH
    SetGas(21000).                    // Standard transfer gas limit
    SetGasPriceInGwei("20").          // 20 Gwei gas price
    SetNonce(42).                     // Transaction nonce
    SetChainID(big.NewInt(1))         // Mainnet chain ID

// Sign the transaction
privateKey, err := web3.PrivateKeyFromHex("0x...")
if err != nil {
    log.Fatal(err)
}

signedTx, err := web3.SignTransaction(txParams, privateKey)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Transaction Hash: %s\n", signedTx.Hash)
fmt.Printf("Raw Transaction: %s\n", signedTx.Raw)
```

#### EIP-1559 Transaction (Type 2)
```go
// Create EIP-1559 transaction parameters
eip1559Params := web3.NewEIP1559TransactionParams()
eip1559Params.To = "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
eip1559Params.Value, _ = web3.ToWei("0.05", "ether")
eip1559Params.Gas = 21000
eip1559Params.MaxFeePerGas, _ = web3.ToWei("30", "gwei")         // Maximum fee per gas
eip1559Params.MaxPriorityFeePerGas, _ = web3.ToWei("2", "gwei")  // Tip to miners
eip1559Params.Nonce = 43
eip1559Params.ChainID = big.NewInt(1)

// Sign EIP-1559 transaction
signedTx, err := web3.SignEIP1559Transaction(eip1559Params, privateKey)
if err != nil {
    log.Fatal(err)
}
```

### High-Level Transaction Methods

#### Send Ether
```go
// Send ETH using wallet (handles nonce, gas estimation automatically)
result, err := wallet.SendEther(ctx, "0xRECIPIENT_ADDRESS", "1.5")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Transaction sent: %s\n", result.TransactionHash)
fmt.Printf("From: %s\n", result.From)
fmt.Printf("To: %s\n", result.To)
```

#### Send with Custom Options
```go
// Send with custom gas settings
transferOpts := &web3.TransferOptions{
    To:       "0xRECIPIENT_ADDRESS",
    Value:    web3.ToWei("2.0", "ether"),
    GasLimit: 25000,                           // Custom gas limit
    GasPrice: web3.ToWei("25", "gwei"),        // Custom gas price
    Data:     []byte("Hello Ethereum!"),       // Optional data
}

result, err := wallet.SendTransaction(ctx, transferOpts)
if err != nil {
    log.Fatal(err)
}
```

#### Send EIP-1559 Transaction
```go
maxFee, _ := web3.ToWei("50", "gwei")
priorityFee, _ := web3.ToWei("3", "gwei")

result, err := wallet.SendEIP1559Transaction(ctx, transferOpts, maxFee, priorityFee)
if err != nil {
    log.Fatal(err)
}
```

### Smart Contract Interactions

#### Contract Method Calls (Read-Only)
```go
// Call a contract method (no transaction, no gas cost)
contractAddress := "0xA0b86a33E6417c48cd7a94Ca95e70aD2c51e74f7"

// Example: balanceOf(address) call
balanceOfData, err := web3.EncodeABI("balanceOf(address)", wallet.GetAddress())
if err != nil {
    log.Fatal(err)
}

result, err := wallet.CallContract(ctx, contractAddress, balanceOfData)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Token balance (hex): %s\n", result)
```

#### Contract Transactions (State-Changing)
```go
// Example: ERC-20 transfer
recipientAddress := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
transferAmount := web3.ToWei("100", "ether") // 100 tokens (assuming 18 decimals)

// Encode transfer(address,uint256) function call
transferData, err := web3.EncodeABI("transfer(address,uint256)", recipientAddress, transferAmount)
if err != nil {
    log.Fatal(err)
}

// Send contract transaction
result, err := wallet.SendContractTransaction(ctx, contractAddress, transferData, big.NewInt(0))
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Token transfer transaction: %s\n", result.TransactionHash)
```

#### Contract Deployment
```go
// Deploy a smart contract
contractBytecode := []byte{0x60, 0x80, 0x60, 0x40, /* ... contract bytecode ... */}
constructorArgs := []byte{/* encoded constructor parameters */}

result, err := wallet.DeployContract(ctx, contractBytecode, constructorArgs, 500000, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Contract deployment transaction: %s\n", result.TransactionHash)

// Wait for deployment confirmation
receipt, err := wallet.WaitForTransaction(ctx, result.TransactionHash)
if err != nil {
    log.Fatal(err)
}

if receipt.ContractAddress != "" {
    fmt.Printf("Contract deployed at: %s\n", receipt.ContractAddress)
}
```

### Transaction Monitoring

#### Wait for Transaction Confirmation
```go
// Wait for transaction to be mined
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

receipt, err := wallet.WaitForTransaction(ctx, txHash)
if err != nil {
    log.Fatal(err)
}

// Check transaction status
if receipt.Status == "0x1" {
    fmt.Println("✅ Transaction successful!")
    fmt.Printf("Gas used: %s\n", receipt.GasUsed)
    fmt.Printf("Block number: %s\n", receipt.BlockNumber)
} else {
    fmt.Println("❌ Transaction failed!")
}
```

### Advanced Transaction Features

#### Batch Operations
```go
// Send multiple transactions concurrently
addresses := []string{
    "0xRecipient1...",
    "0xRecipient2...",
    "0xRecipient3...",
}

type TxResult struct {
    Address string
    TxHash  string
    Error   error
}

results := make(chan TxResult, len(addresses))

for _, addr := range addresses {
    go func(recipient string) {
        result, err := wallet.SendEther(ctx, recipient, "0.01")
        if err != nil {
            results <- TxResult{recipient, "", err}
            return
        }
        results <- TxResult{recipient, result.TransactionHash, nil}
    }(addr)
}

// Collect results
for i := 0; i < len(addresses); i++ {
    result := <-results
    if result.Error != nil {
        fmt.Printf("❌ Failed to send to %s: %v\n", result.Address, result.Error)
    } else {
        fmt.Printf("✅ Sent to %s: %s\n", result.Address, result.TxHash)
    }
}
```

#### Gas Estimation and Optimization
```go
// Get optimal gas price with buffer
currentGas, err := client.Eth().GetGasPrice(ctx)
if err != nil {
    log.Fatal(err)
}

// Add 10% buffer for faster confirmation
buffer := new(big.Int).Div(currentGas, big.NewInt(10))
optimalGas := new(big.Int).Add(currentGas, buffer)

fmt.Printf("Optimal gas price: %s Gwei\n", 
    web3.FromWei(optimalGas, "gwei"))

// Use in transaction
txParams := web3.NewTransactionParams().
    SetTo("0x...").
    SetValueInEther("1.0").
    SetGasPrice(optimalGas)
```

#### Transaction Recovery
```go
// Recover the signer address from a raw transaction
rawTxHex := "0xf86c808504a817c800825208940x..."
signerAddress, err := web3.RecoverSigner(rawTxHex)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Transaction was signed by: %s\n", signerAddress)
```

## 🎯 Typed Constants & Enums

The library eliminates hardcoded strings with comprehensive typed constants for better type safety and IDE support.

### Block Parameters

Use typed block parameters instead of strings:

```go
// Instead of hardcoded strings
balance, err := client.Eth().GetBalance(ctx, address, "latest")

// Use typed constants
balance, err := client.Eth().GetBalance(ctx, address, web3.BlockLatest)
balance, err := client.Eth().GetBalance(ctx, address, web3.BlockPending)
balance, err := client.Eth().GetBalance(ctx, address, web3.BlockEarliest)

// For specific block numbers
blockParam := web3.BlockNumber(18500000)
balance, err := client.Eth().GetBalance(ctx, address, blockParam)
```

### Ether Units

Type-safe unit conversions:

```go
// Instead of string units
weiValue, _ := web3.ToWei("1", "ether")
gweiValue, _ := web3.ToWei("20", "gwei")

// Use typed units
weiValue, _ := web3.ToWei("1", web3.Ether)
gweiValue, _ := web3.ToWei("20", web3.Gwei)

// Helper functions
ethWei, _ := web3.EtherToWei("2.5")
gweiWei, _ := web3.GweiToWei("30")

ethDisplay, _ := web3.WeiToEther(ethWei)
gweiDisplay, _ := web3.WeiToGwei(gweiWei)
```

Available units: `Wei`, `Kwei`, `Mwei`, `Gwei`, `Szabo`, `Finney`, `Ether`, `Kether`, `Mether`, `Gether`, `Tether`

### Chain IDs

Pre-defined chain IDs for popular networks:

```go
// Instead of magic numbers
txParams.SetChainID(big.NewInt(1))

// Use typed chain IDs
txParams.SetChainID(web3.ChainMainnet)
txParams.SetChainID(web3.ChainGoerli)
txParams.SetChainID(web3.ChainPolygon)
txParams.SetChainID(web3.ChainArbitrum)

// Get network information
config, err := web3.GetNetworkConfig(web3.ChainMainnet)
fmt.Printf("Network: %s, Currency: %s\n", config.Name, config.Currency)

// Check network type
isTestnet := web3.IsTestnet(web3.ChainGoerli) // true
isMainnet := web3.IsMainnet(web3.ChainMainnet) // true
```

### Gas Limits

Standard gas limits for common operations:

```go
// Predefined gas limits
txParams.SetGas(web3.GasLimitTransfer.Uint64())        // 21,000
txParams.SetGas(web3.GasLimitTokenTransfer.Uint64())   // 65,000
txParams.SetGas(web3.GasLimitTokenApproval.Uint64())   // 50,000
txParams.SetGas(web3.GasLimitContractDeploy.Uint64())  // 500,000
```

### Gas Price Optimization

Smart gas pricing with predefined levels:

```go
// Gas price optimization levels
optimal, err := web3.GetOptimalGasPrice(ctx, client, web3.GasPriceStandard) // +10%
rapid, err := web3.GetOptimalGasPrice(ctx, client, web3.GasPriceRapid)     // +50%

// Levels: GasPriceSlow, GasPriceStandard, GasPriceFast, GasPriceRapid
```

### Common Addresses

Pre-defined addresses for popular contracts:

```go
// Common Ethereum addresses
wethContract := web3.WETHMainnet.String()
usdcContract := web3.USDCMainnet.String()
usdtContract := web3.USDTMainnet.String()

// Special addresses
zeroAddr := web3.ZeroAddress.String()
burnAddr := web3.BurnAddress.String()

// Address validation helpers
if web3.IsZeroAddress(someAddress) {
    fmt.Println("Cannot send to zero address")
}
```

### Function Signatures

Type-safe ERC-20 function signatures:

```go
// Instead of hardcoded strings
data, err := web3.EncodeABI("balanceOf(address)", address)

// Use typed signatures
data, err := web3.EncodeABI(web3.FuncBalanceOf.String(), address)
data, err := web3.EncodeABI(web3.FuncTransfer.String(), to, amount)
data, err := web3.EncodeABI(web3.FuncApprove.String(), spender, amount)
```

### Transaction Helpers

High-level transaction builders:

```go
// Simple ETH transfer
ethTx := web3.NewSimpleTransfer("0xRecipient", "1.5", web3.ChainMainnet)

// Token transfer
tokenTx, err := web3.NewTokenTransfer(
    web3.USDCMainnet.String(),
    "0xRecipient", 
    amount,
    web3.ChainMainnet,
)

// Token approval
approvalTx, err := web3.NewTokenApproval(
    web3.USDCMainnet.String(),
    spenderAddress,
    amount,
    web3.ChainMainnet,
)
```

### Transaction Status

Type-safe transaction status checking:

```go
// Instead of string comparisons
if receipt.Status == "0x1" {
    fmt.Println("Success")
}

// Use typed status
if web3.IsTransactionSuccess(receipt) {
    fmt.Println("Transaction successful!")
}

if web3.IsTransactionFailure(receipt) {
    fmt.Println("Transaction failed!")
}
```

### Complete Example

```go
// Before: Hardcoded strings everywhere
balance, _ := client.Eth().GetBalance(ctx, addr, "latest")
weiVal, _ := web3.ToWei("1.5", "ether")
gasPrice, _ := web3.ToWei("20", "gwei")

tx := web3.NewTransactionParams().
    SetTo(recipient).
    SetValue(weiVal).
    SetGas(21000).
    SetGasPrice(gasPrice).
    SetChainID(big.NewInt(1))

// After: Type-safe constants
balance, _ := client.Eth().GetBalance(ctx, addr, web3.BlockLatest)
weiVal, _ := web3.ToWei("1.5", web3.Ether)
gasPrice, _ := web3.ToWei("20", web3.Gwei)

tx := web3.NewTransactionParams().
    SetTo(recipient).
    SetValue(weiVal).
    SetGas(web3.GasLimitTransfer.Uint64()).
    SetGasPrice(gasPrice).
    SetChainID(web3.ChainMainnet)
```

### Benefits

- **Type Safety**: Compile-time checking prevents typos
- **IDE Support**: Autocomplete and IntelliSense 
- **Maintainability**: Centralized constant definitions
- **Readability**: Self-documenting code
- **Consistency**: Standardized values across the library

## 🛠️ Utility Functions

### Wei/Ether Conversion

The library provides comprehensive unit conversion similar to web3.js:

#### Convert to Wei
```go
// Convert from Ether to Wei
weiValue, err := web3.ToWei("1.5", "ether")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("1.5 ETH = %s Wei\n", weiValue.String())

// Convert from Gwei to Wei  
gasPriceWei, err := web3.ToWei("20", "gwei")
fmt.Printf("20 Gwei = %s Wei\n", gasPriceWei.String())

// Support for all units
units := []string{"wei", "kwei", "mwei", "gwei", "szabo", "finney", "ether"}
for _, unit := range units {
    wei, _ := web3.ToWei("1", unit)
    fmt.Printf("1 %s = %s wei\n", unit, wei.String())
}
```

#### Convert from Wei
```go
// Convert Wei to Ether for display
weiAmount := big.NewInt(1500000000000000000) // 1.5 ETH in wei
ethValue, err := web3.FromWei(weiAmount, "ether")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Wei to ETH: %s\n", ethValue) // "1.5"

// Convert Wei to Gwei (useful for gas prices)
gasPriceWei := big.NewInt(20000000000)
gasPriceGwei, _ := web3.FromWei(gasPriceWei, "gwei")
fmt.Printf("Gas price: %s Gwei\n", gasPriceGwei) // "20"
```

### Address Validation

#### Check Address Format
```go
addresses := []string{
    "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // valid
    "0xInvalidAddress",                            // invalid
    "d8dA6BF26964aF9D7eEd9e03E53415D37aA96045",   // missing 0x prefix
    "0x123",                                       // too short
}

for _, addr := range addresses {
    if web3.IsAddress(addr) {
        fmt.Printf("%s is a valid address\n", addr)
    } else {
        fmt.Printf("%s is invalid\n", addr)
    }
}
```

### Hex Conversion

#### Convert Values to Hex
```go
// Convert integers to hex
hexInt := web3.ToHex(12345)        // "0x3039"
hexBigInt := web3.ToHex(big.NewInt(12345)) // "0x3039"

// Convert strings to hex
hexString := web3.ToHex("hello")   // "0x68656c6c6f"

// Convert byte arrays to hex
data := []byte{0x12, 0x34, 0x56}
hexBytes := web3.ToHex(data)       // "0x123456"
```

#### Convert from Hex
```go
// Convert hex string to big.Int
hexValue := "0x3039"
value, err := web3.FromHex(hexValue)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Hex %s = %d decimal\n", hexValue, value.Int64()) // "12345"
```

### String Padding

#### Pad Strings
```go
// Left padding (useful for hex values)
padded := web3.PadLeft("abc", 8, "0")    // "00000abc"
hexPadded := web3.PadLeft("1a2b", 8, "0") // "00001a2b"

// Right padding
rightPadded := web3.PadRight("abc", 8, "0") // "abc00000"

// Common use case: padding addresses or hex values
address := "0x123"
paddedAddr := "0x" + web3.PadLeft(address[2:], 40, "0")
fmt.Printf("Padded address: %s\n", paddedAddr)
```

## 📏 Supported Ethereum Units

| Unit Name | Aliases | Wei Value | Common Use |
|-----------|---------|-----------|------------|
| `wei` | - | 1 | Smallest unit, precise calculations |
| `kwei` | `babbage`, `femtoether` | 10³ | - |
| `mwei` | `lovelace`, `picoether` | 10⁶ | - |
| `gwei` | `shannon`, `nanoether`, `nano` | 10⁹ | **Gas prices** |
| `szabo` | `microether`, `micro` | 10¹² | - |
| `finney` | `milliether`, `milli` | 10¹⁵ | Small transactions |
| `ether` | `eth` | 10¹⁸ | **Standard currency unit** |
| `kether` | `grand` | 10²¹ | Large amounts |
| `mether` | - | 10²⁴ | Very large amounts |
| `gether` | - | 10²⁷ | Extremely large amounts |
| `tether` | - | 10³⁰ | Theoretical amounts |

### Unit Conversion Examples
```go
// Common conversions
oneEth, _ := web3.ToWei("1", "ether")           // 1000000000000000000 wei
oneGwei, _ := web3.ToWei("1", "gwei")           // 1000000000 wei
gasPrice, _ := web3.ToWei("20", "gwei")         // 20000000000 wei (20 Gwei)

// Display conversions
weiAmount := big.NewInt(1500000000000000000)
ethDisplay, _ := web3.FromWei(weiAmount, "ether")    // "1.5"
gweiDisplay, _ := web3.FromWei(weiAmount, "gwei")    // "1500000000"
```

## 🚨 Error Handling

### RPC Errors
The library provides detailed error information for RPC failures:

```go
balance, err := client.Eth().GetBalance(ctx, "invalid-address", "latest")
if err != nil {
    // Check if it's an RPC error
    if rpcErr, ok := err.(*web3.RPCError); ok {
        fmt.Printf("RPC Error %d: %s\n", rpcErr.Code, rpcErr.Message)
        if rpcErr.Data != "" {
            fmt.Printf("Additional data: %s\n", rpcErr.Data)
        }
    } else {
        // Handle other errors (network, parsing, etc.)
        fmt.Printf("Other error: %v\n", err)
    }
}
```

### Common Error Scenarios
```go
// Timeout handling
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := client.Eth().GetBlockNumber(ctx)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        fmt.Println("Request timed out")
    }
    return
}

// Network error handling
if err != nil {
    if strings.Contains(err.Error(), "connection refused") {
        fmt.Println("Cannot connect to Ethereum node")
    } else if strings.Contains(err.Error(), "invalid response") {
        fmt.Println("Invalid response from node")
    }
}
```

## 🌟 Advanced Usage Examples

### Complete Transaction Monitoring
```go
func monitorTransaction(client *web3.Client, txHash string) error {
    ctx := context.Background()
    
    // Check if transaction exists
    tx, err := client.Eth().GetTransactionByHash(ctx, txHash)
    if err != nil {
        return fmt.Errorf("transaction not found: %w", err)
    }
    
    fmt.Printf("Transaction found: %s -> %s\n", tx.From, tx.To)
    fmt.Printf("Value: %s ETH\n", mustFromWei(tx.Value, "ether"))
    
    // Wait for confirmation
    for {
        receipt, err := client.Eth().GetTransactionReceipt(ctx, txHash)
        if err != nil {
            time.Sleep(5 * time.Second)
            continue
        }
        
        if receipt.Status == "0x1" {
            fmt.Printf("✅ Transaction confirmed in block %s\n", receipt.BlockNumber)
            fmt.Printf("Gas used: %s\n", receipt.GasUsed)
            return nil
        } else {
            fmt.Printf("❌ Transaction failed\n")
            return fmt.Errorf("transaction failed")
        }
    }
}

func mustFromWei(weiHex, unit string) string {
    wei, _ := web3.FromHex(weiHex)
    eth, _ := web3.FromWei(wei, unit)
    return eth
}
```

### Gas Price Optimization
```go
func getOptimalGasPrice(client *web3.Client) (*big.Int, error) {
    ctx := context.Background()
    
    // Get current gas price
    currentGas, err := client.Eth().GetGasPrice(ctx)
    if err != nil {
        return nil, err
    }
    
    // Add 10% buffer for faster processing
    buffer := new(big.Int).Div(currentGas, big.NewInt(10))
    optimalGas := new(big.Int).Add(currentGas, buffer)
    
    fmt.Printf("Current gas price: %s Gwei\n", mustFromWei(fmt.Sprintf("0x%x", currentGas), "gwei"))
    fmt.Printf("Optimal gas price: %s Gwei\n", mustFromWei(fmt.Sprintf("0x%x", optimalGas), "gwei"))
    
    return optimalGas, nil
}
```

### Batch Operations
```go
func getMultipleBalances(client *web3.Client, addresses []string) (map[string]*big.Int, error) {
    ctx := context.Background()
    balances := make(map[string]*big.Int)
    
    // Use goroutines for concurrent requests
    type result struct {
        address string
        balance *big.Int
        err     error
    }
    
    results := make(chan result, len(addresses))
    
    for _, addr := range addresses {
        go func(address string) {
            balance, err := client.Eth().GetBalance(ctx, address, "latest")
            results <- result{address, balance, err}
        }(addr)
    }
    
    // Collect results
    for i := 0; i < len(addresses); i++ {
        res := <-results
        if res.err != nil {
            return nil, res.err
        }
        balances[res.address] = res.balance
    }
    
    return balances, nil
}
```

## 📁 Project Structure

```
go-web3/
├── client.go          # Core RPC client implementation
├── eth.go             # Ethereum-specific methods
├── utils.go           # Utility functions  
├── example/
│   └── main.go        # Usage examples
├── go.mod             # Go module definition
├── README.md          # This documentation
└── LICENSE            # MIT license
```

## 🔗 Popular RPC Endpoints

### Mainnet
- **Infura**: `https://mainnet.infura.io/v3/YOUR_PROJECT_ID`
- **Alchemy**: `https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY`
- **QuickNode**: `https://your-endpoint.quiknode.pro/YOUR_API_KEY/`

### Testnets  
- **Goerli**: `https://goerli.infura.io/v3/YOUR_PROJECT_ID`
- **Sepolia**: `https://sepolia.infura.io/v3/YOUR_PROJECT_ID`

### Local Development
- **Hardhat**: `http://localhost:8545`
- **Ganache**: `http://localhost:7545`

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Support

- 📖 Documentation: This README
- 💻 Examples: See `example/main.go`
- 🐛 Issues: Open an issue on GitHub
- 💡 Feature Requests: Open an issue with the enhancement label

---

**Made with ❤️ for the Ethereum Go community**