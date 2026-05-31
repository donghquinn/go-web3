# Go-Web3
[![Test](https://github.com/donghquinn/go-web3/actions/workflows/test.yml/badge.svg)](https://github.com/donghquinn/go-web3/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/donghquinn/go-web3/branch/main/graph/badge.svg)](https://codecov.io/gh/donghquinn/go-web3)

A comprehensive Go library for interacting with the Ethereum blockchain through JSON-RPC APIs, designed with a web3.js-like interface for familiar usage patterns.

## Features

- **Complete Ethereum JSON-RPC**: All major `eth_*`, `net_*`, and `web3_*` methods
- **EIP-1559 Support**: Fee history, max priority fee, and type-2 transaction signing
- **EIP-4844 Awareness**: Blob gas fields included in fee history responses
- **Event Log Queries**: Filter contract events with `eth_getLogs`
- **Contract Inspection**: Read bytecode and storage slots at any block
- **Pending Transaction Support**: Monitor mempool and get pending transactions by account
- **Context-Aware Operations**: Proper context handling for timeouts and cancellation
- **Type-Safe Structures**: Strongly-typed transaction, block, receipt, log, and fee objects
- **Built-in Utilities**: Wei/Ether conversion, address validation, hex operations
- **Robust Error Handling**: Detailed RPC error information with proper Go error wrapping
- **Production Ready**: Thread-safe client with connection pooling
- **High Test Coverage**: 97%+ test coverage via `httptest`-based mock RPC servers

## Installation

```bash
go get github.com/donghquinn/go-web3
```

## Quick Start

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
    client := web3.NewClient("https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    blockNumber, err := client.Eth().GetBlockNumber(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Latest block: %d\n", blockNumber)

    address := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
    balance, err := client.Eth().GetBalance(ctx, address, web3.BlockLatest)
    if err != nil {
        log.Fatal(err)
    }

    balanceEth, _ := web3.FromWei(balance, web3.Ether)
    fmt.Printf("Balance: %s ETH\n", balanceEth)
}
```

## API Reference

### Creating a Client

```go
client := web3.NewClient("YOUR_ETHEREUM_RPC_URL")

// Popular RPC endpoints:
// Mainnet: "https://mainnet.infura.io/v3/YOUR_PROJECT_ID"
// Alchemy: "https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY"
// Local:   "http://localhost:8545"
```

### Account & Balance Operations

```go
// Get balance in Wei
balance, err := client.Eth().GetBalance(ctx, address, web3.BlockLatest)

// Convert to Ether for display
balanceEth, _ := web3.FromWei(balance, web3.Ether)
fmt.Printf("Balance: %s ETH\n", balanceEth)

// At a specific block
balance, err = client.Eth().GetBalance(ctx, address, web3.BlockNumber(18000000))

// Get nonce (transaction count)
nonce, err := client.Eth().GetTransactionCount(ctx, address, web3.BlockLatest)

// Pending nonce (for rapid transaction sending)
pendingNonce, err := client.Eth().GetTransactionCount(ctx, address, web3.BlockPending)
```

### Block Operations

```go
// Latest block number
blockNumber, err := client.Eth().GetBlockNumber(ctx)

// Block by number (false = no full tx bodies)
block, err := client.Eth().GetBlockByNumber(ctx, web3.BlockLatest, false)
fmt.Printf("Hash: %s, GasUsed: %s\n", block.Hash, block.GasUsed)

// Block with full transaction details
block, err = client.Eth().GetBlockByNumber(ctx, web3.BlockLatest, true)

// Block by hash
block, err = client.Eth().GetBlockByHash(ctx, "0xabc...", false)
```

### Transaction Operations

```go
// Get transaction by hash
tx, err := client.Eth().GetTransactionByHash(ctx, "0xabc...")
fmt.Printf("From: %s, To: %s, Value: %s\n", tx.From, tx.To, tx.Value)

// Get receipt
receipt, err := client.Eth().GetTransactionReceipt(ctx, "0xabc...")
if receipt.Status == "0x1" {
    fmt.Println("Transaction successful")
}

// Send pre-signed raw transaction
txHash, err := client.Eth().SendRawTransaction(ctx, "0xf86c...")
```

### Pending Transaction Operations

```go
// All pending transactions in mempool
pendingTxs, err := client.Eth().GetPendingTransactions(ctx)

// Pending transaction count
count, err := client.Eth().GetPendingTransactionCount(ctx)

// Pending transactions for a specific account (from or to)
accountTxs, err := client.Eth().GetAccountPendingTransactions(ctx, address)

// Check if a transaction is pending
isPending, err := client.Eth().IsPendingTransaction(ctx, txHash)
```

### Gas Operations

```go
// Current gas price
gasPrice, err := client.Eth().GetGasPrice(ctx)
gasPriceGwei, _ := web3.FromWei(gasPrice, web3.Gwei)
fmt.Printf("Gas price: %s Gwei\n", gasPriceGwei)

// EIP-1559: suggested tip
tip, err := client.Eth().GetMaxPriorityFeePerGas(ctx)

// Fee history for last 10 blocks
feeHistory, err := client.Eth().GetFeeHistory(ctx, 10, web3.BlockLatest, []float64{25, 75})
for i, baseFee := range feeHistory.BaseFeePerGas {
    fmt.Printf("Block +%d base fee: %s\n", i, baseFee)
}

// Estimate gas for a transaction
gasEstimate, err := client.Eth().EstimateGas(ctx, map[string]interface{}{
    "from":  "0xSENDER",
    "to":    "0xRECIPIENT",
    "value": "0xde0b6b3a7640000", // 1 ETH
})
```

### Contract Interaction

```go
// Read-only contract call
callObj := map[string]interface{}{
    "to":   "0xCONTRACT",
    "data": "0x70a08231000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045",
}
result, err := client.Eth().Call(ctx, callObj, web3.BlockLatest)

// Get contract bytecode
code, err := client.Eth().GetCode(ctx, "0xCONTRACT", web3.BlockLatest)
if code == "0x" {
    fmt.Println("Not a contract")
}

// Read storage slot
slot, err := client.Eth().GetStorageAt(ctx, "0xCONTRACT", "0x0", web3.BlockLatest)
```

### Event Logs

```go
// Filter Transfer events
filter := web3.LogFilter{
    Address:   "0xCONTRACT",
    FromBlock: web3.BlockNumber(18000000),
    ToBlock:   web3.BlockLatest,
    Topics: []interface{}{
        "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", // Transfer
    },
}
logs, err := client.Eth().GetLogs(ctx, filter)

// Multi-address, OR-topic filter
filter = web3.LogFilter{
    Address: []string{"0xContractA", "0xContractB"},
    Topics: []interface{}{
        []string{
            "0xddf252ad...", // Transfer
            "0x8c5be1e5...", // Approval
        },
    },
    FromBlock: web3.BlockLatest,
    ToBlock:   web3.BlockLatest,
}
```

### Network Info

```go
chainID, err := client.Eth().GetChainID(ctx)
fmt.Printf("Chain ID: %s\n", chainID.String()) // "1" for mainnet

netVersion, err := client.Eth().GetNetVersion(ctx)

clientVersion, err := client.Eth().GetClientVersion(ctx)
fmt.Printf("Client: %s\n", clientVersion) // e.g. "Geth/v1.13.0/..."
```

---

## Transaction Management

### Wallet

```go
// Load wallet from private key
wallet, err := web3.NewWallet("0xPRIVATE_KEY", client)

// Generate a new random wallet
wallet, err = web3.CreateWallet(client)
fmt.Printf("Address:     %s\n", wallet.GetAddress())
fmt.Printf("Private Key: %s\n", wallet.GetPrivateKey()) // keep secure!

// Check balance
balance, err := wallet.GetBalance(ctx)

// Get nonce
nonce, err := wallet.GetNonce(ctx)
```

### Sending ETH

```go
// Simple ETH transfer (gas estimated automatically)
result, err := wallet.SendEther(ctx, "0xRECIPIENT", "1.5")
fmt.Printf("TxHash: %s\n", result.TransactionHash)

// Send wei directly
result, err = wallet.SendWei(ctx, "0xRECIPIENT", big.NewInt(1000))

// Custom gas settings
result, err = wallet.SendTransaction(ctx, &web3.TransferOptions{
    To:       "0xRECIPIENT",
    Value:    big.NewInt(1e18),
    GasLimit: 25000,
    GasPrice: big.NewInt(25e9),
})

// EIP-1559 transaction
maxFee, _ := web3.ToWei("50", web3.Gwei)
tip, _ := web3.ToWei("3", web3.Gwei)
result, err = wallet.SendEIP1559Transaction(ctx, &web3.TransferOptions{
    To:    "0xRECIPIENT",
    Value: big.NewInt(0),
}, maxFee, tip)
```

### Smart Contracts

```go
// Read-only call via wallet
balanceOfData, err := web3.EncodeABI("balanceOf(address)", wallet.GetAddress())
result, err := wallet.CallContract(ctx, "0xCONTRACT", balanceOfData)

// State-changing transaction
transferData, err := web3.EncodeABI("transfer(address,uint256)",
    "0xRECIPIENT",
    big.NewInt(100e18),
)
result, err = wallet.SendContractTransaction(ctx, "0xCONTRACT", transferData, big.NewInt(0))

// Contract deployment (gasLimit=0 triggers estimation)
result, err = wallet.DeployContract(ctx, bytecode, constructorArgs, 0, nil)
fmt.Printf("Deploy tx: %s\n", result.TransactionHash)

// Wait for receipt
receipt, err := wallet.WaitForTransaction(ctx, result.TransactionHash)
if receipt.ContractAddress != "" {
    fmt.Printf("Deployed at: %s\n", receipt.ContractAddress)
}
```

### Building and Signing Transactions Manually

```go
// Legacy (pre-EIP-1559)
txParams := web3.NewTransactionParams().
    SetTo("0xRECIPIENT").
    SetValueInEther("0.1").
    SetGas(web3.GasLimitTransfer.Uint64()).
    SetGasPriceInGwei("20").
    SetNonce(42).
    SetChainID(web3.ChainMainnet)

privateKey, err := web3.PrivateKeyFromHex("0xPRIVATE_KEY")
signedTx, err := web3.SignTransaction(txParams, privateKey)
fmt.Printf("Hash: %s\nRaw:  %s\n", signedTx.Hash, signedTx.Raw)

// EIP-1559 (type 2)
eip1559 := web3.NewEIP1559TransactionParams()
eip1559.To = "0xRECIPIENT"
eip1559.Gas = 21000
eip1559.MaxFeePerGas, _ = web3.ToWei("30", web3.Gwei)
eip1559.MaxPriorityFeePerGas, _ = web3.ToWei("2", web3.Gwei)
eip1559.Nonce = 43
eip1559.ChainID = web3.ChainMainnet.BigInt()

signedTx, err = web3.SignEIP1559Transaction(eip1559, privateKey)

// Contract deployment (empty To = contract creation)
signed, err := web3.CreateContractDeployment(bytecode, constructorData, privateKey, txParams)

// Contract call
signed, err = web3.CreateContractCall("0xCONTRACT", methodData, privateKey, txParams)

// Recover signer from raw tx
signer, err := web3.RecoverSigner("0xf86c...")
```

---

## Typed Constants

### Block Parameters

```go
web3.BlockLatest    // "latest"
web3.BlockPending   // "pending"
web3.BlockEarliest  // "earliest"
web3.BlockNumber(18500000)  // specific block
```

### Ether Units

```go
weiValue, _ := web3.ToWei("1.5", web3.Ether)
gweiValue, _ := web3.ToWei("20", web3.Gwei)

// Convenience helpers
ethWei, _ := web3.EtherToWei("2.5")
gweiWei, _ := web3.GweiToWei("30")

ethDisplay, _ := web3.WeiToEther(ethWei)
gweiDisplay, _ := web3.WeiToGwei(gweiWei)
```

| Unit | Aliases | Wei Value |
|------|---------|-----------|
| `Wei` | — | 10⁰ |
| `Kwei` | `Babbage`, `Femtoether` | 10³ |
| `Mwei` | `Lovelace`, `Picoether` | 10⁶ |
| `Gwei` | `Shannon`, `Nanoether`, `Nano` | 10⁹ |
| `Szabo` | `Microether`, `Micro` | 10¹² |
| `Finney` | `Milliether`, `Milli` | 10¹⁵ |
| `Ether` | `EthUnit` | 10¹⁸ |
| `Kether` | `Grand` | 10²¹ |
| `Mether` | — | 10²⁴ |
| `Gether` | — | 10²⁷ |
| `Tether` | — | 10³⁰ |

### Chain IDs

```go
web3.ChainMainnet   // 1
web3.ChainGoerli    // 5
web3.ChainSepolia   // 11155111
web3.ChainPolygon   // 137
web3.ChainArbitrum  // 42161

config, err := web3.GetNetworkConfig(web3.ChainMainnet)
fmt.Printf("Network: %s, Currency: %s\n", config.Name, config.Currency)

web3.IsTestnet(web3.ChainGoerli)  // true
web3.IsMainnet(web3.ChainMainnet) // true
```

### Gas Limits

```go
web3.GasLimitTransfer.Uint64()       // 21,000
web3.GasLimitTokenTransfer.Uint64()  // 65,000
web3.GasLimitTokenApproval.Uint64()  // 50,000
web3.GasLimitContractDeploy.Uint64() // 500,000
```

### Gas Price Levels

```go
price, err := web3.GetOptimalGasPrice(ctx, client, web3.GasPriceStandard) // base + 10%
price, err  = web3.GetOptimalGasPrice(ctx, client, web3.GasPriceFast)     // base + 25%
price, err  = web3.GetOptimalGasPrice(ctx, client, web3.GasPriceRapid)    // base + 50%
// Also: GasPriceSlow (base + 0%)
```

### Common Addresses

```go
web3.WETHMainnet.String()  // Wrapped ETH contract
web3.USDCMainnet.String()  // USDC contract
web3.USDTMainnet.String()  // USDT contract
web3.ZeroAddress.String()  // 0x0000...0000
web3.BurnAddress.String()  // 0xdead...dead

web3.IsZeroAddress(addr)
web3.IsBurnAddress(addr)
```

### Function Signatures

```go
data, err := web3.EncodeABI(web3.FuncBalanceOf.String(), address)
data, err  = web3.EncodeABI(web3.FuncTransfer.String(), to, amount)
data, err  = web3.EncodeABI(web3.FuncApprove.String(), spender, amount)
```

### Transaction Status Helpers

```go
if web3.IsTransactionSuccess(receipt) {
    fmt.Println("Transaction successful!")
}
if web3.IsTransactionFailure(receipt) {
    fmt.Println("Transaction failed!")
}

fee := web3.CalculateTransactionFee(gasLimit, gasPrice)
```

---

## go-blockchain-helper Integration

This library integrates [`go-blockchain-helper`](https://github.com/donghquinn/go-blockchain-helper) for advanced ABI encoding and token utilities.

### ERC-20 Token Support

```go
// Create a typed token instance
token := web3.NewERC20Token("0xCONTRACT", "USD Coin", "USDC", 6)

// Encode common operations
transferData, err := web3.EncodeERC20Transfer(token, recipient, amount)
approveData, err  := web3.EncodeERC20Approve(token, spender, amount)
fromData, err     := web3.EncodeERC20TransferFrom(token, from, to, amount)

// High-level transaction builders
tokenTx, err := web3.NewTokenTransfer(
    web3.USDCMainnet.String(), recipient, amount, web3.ChainMainnet,
)
approvalTx, err := web3.NewTokenApproval(
    web3.USDCMainnet.String(), spender, amount, web3.ChainMainnet,
)

// Query on-chain values
balance, err   := web3.GetTokenBalance(ctx, client, "0xCONTRACT", ownerAddr)
allowance, err := web3.GetTokenAllowance(ctx, client, "0xCONTRACT", owner, spender)
```

### ERC-721 (NFT) Support

```go
nft := web3.NewERC721Token("0xCONTRACT", "BoredApeYachtClub", "BAYC")

transferData, err     := web3.EncodeERC721Transfer(nft, from, to, tokenId)
approveData, err      := web3.EncodeERC721Approve(nft, approved, tokenId)
setApprovalData, err  := web3.EncodeERC721SetApprovalForAll(nft, operator, true)
```

### Advanced ABI Encoding

```go
// Simple encoding with type inference
data, err := web3.EncodeABI("transfer(address,uint256)", recipientAddr, amount)

// Advanced encoding with explicit type hints
params := []blockchainhelper.ABIParam{
    {Type: "address"},
    {Type: "uint256"},
}
data, err = web3.EncodeFunctionCallAdvanced("transfer(address,uint256)", params, args)

// Decode results
values, err := web3.DecodeFunctionResult([]string{"uint256"}, rawBytes)

// Parse Transfer events
event, err := web3.ParseTransferEvent(logEntry)
```

### Gas Estimation

```go
// Transaction with library-estimated gas
tx, err := web3.CreateTransactionWithEstimate(recipient, value, data, web3.ChainMainnet)

// Manual estimation with buffer
gas, err := web3.EstimateGasWithBuffer(ctx, client, callObj, 0.1) // +10%
```

---

## Utilities

### Address Validation

```go
web3.IsAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e") // true
web3.IsAddress("not-an-address")                             // false

// Private key utilities
addr := web3.PrivateKeyToAddress(privateKey)
hexKey := web3.PrivateKeyToHex(privateKey)
key, err := web3.PrivateKeyFromHex("0xPRIVATE_KEY")

addrStr, err := web3.PrivateKeyToAddressHelper(privateKey)
```

### Hex Conversion

```go
web3.ToHex(12345)                    // "0x3039"
web3.ToHex(big.NewInt(12345))        // "0x3039"
web3.ToHex([]byte{0x12, 0x34})      // "0x1234"
web3.ToHex("hello")                  // "0x68656c6c6f"

val, err := web3.FromHex("0x3039")  // big.Int(12345)
```

### String Padding

```go
web3.PadLeft("abc", 8, "0")   // "00000abc"
web3.PadRight("abc", 8, "0")  // "abc00000"
```

---

## Error Handling

```go
balance, err := client.Eth().GetBalance(ctx, address, web3.BlockLatest)
if err != nil {
    if rpcErr, ok := err.(*web3.RPCError); ok {
        fmt.Printf("RPC error %d: %s\n", rpcErr.Code, rpcErr.Message)
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}

// Timeout handling
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := client.Eth().GetBlockNumber(ctx)
if err != nil && ctx.Err() == context.DeadlineExceeded {
    fmt.Println("Request timed out")
}
```

---

## Project Structure

```
go-web3/
├── client.go           # Core JSON-RPC HTTP client
├── eth.go              # Ethereum API methods (eth_*, net_*, web3_*)
├── wallet.go           # Wallet creation and high-level send helpers
├── transaction.go      # Transaction signing (legacy & EIP-1559)
├── types.go            # Typed constants, enums, and network configs
├── utils.go            # Wei/Ether conversion, hex, address utilities
├── helpers.go          # Advanced helpers (ERC-20/ERC-721, ABI encoding)
├── client_test.go      # JSON-RPC client tests
├── eth_test.go         # Ethereum method tests
├── wallet_test.go      # Wallet and high-level send tests
├── transaction_test.go # Signing and ABI encoding tests
├── helpers_test.go     # Helper function tests
├── utils_test.go       # Utility function tests
├── mock_test.go        # Shared httptest mock RPC server helpers
├── example/
│   └── main.go         # Usage examples
├── go.mod
├── README.md
└── LICENSE
```

## Popular RPC Endpoints

**Mainnet**
- Infura: `https://mainnet.infura.io/v3/YOUR_PROJECT_ID`
- Alchemy: `https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY`
- QuickNode: `https://your-endpoint.quiknode.pro/YOUR_API_KEY/`

**Testnets**
- Sepolia: `https://sepolia.infura.io/v3/YOUR_PROJECT_ID`
- Goerli: `https://goerli.infura.io/v3/YOUR_PROJECT_ID`

**Local Development**
- Hardhat: `http://localhost:8545`
- Ganache: `http://localhost:7545`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes (`git commit -m 'Add my feature'`)
4. Push to the branch (`git push origin feature/my-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
