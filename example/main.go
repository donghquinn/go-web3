package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/donghquinn/go-web3"
)

func basicExample() {
	client := web3.NewClient("https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY")

	ctx := context.Background()

	fmt.Println("=== Ethereum Go-Web3 Library Example ===")

	blockNumber, err := client.Eth().GetBlockNumber(ctx)
	if err != nil {
		log.Printf("Error getting block number: %v", err)
	} else {
		fmt.Printf("Latest block number: %d\n", blockNumber)
	}

	gasPrice, err := client.Eth().GetGasPrice(ctx)
	if err != nil {
		log.Printf("Error getting gas price: %v", err)
	} else {
		gasPriceGwei, _ := web3.FromWei(gasPrice, "gwei")
		fmt.Printf("Current gas price: %s Gwei\n", gasPriceGwei)
	}

	address := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
	balance, err := client.Eth().GetBalance(ctx, address, "latest")
	if err != nil {
		log.Printf("Error getting balance: %v", err)
	} else {
		balanceEth, _ := web3.FromWei(balance, "ether")
		fmt.Printf("Balance of %s: %s ETH\n", address, balanceEth)
	}

	nonce, err := client.Eth().GetTransactionCount(ctx, address, "latest")
	if err != nil {
		log.Printf("Error getting transaction count: %v", err)
	} else {
		fmt.Printf("Transaction count (nonce) for %s: %d\n", address, nonce)
	}

	block, err := client.Eth().GetBlockByNumber(ctx, "latest", false)
	if err != nil {
		log.Printf("Error getting latest block: %v", err)
	} else {
		fmt.Printf("Latest block hash: %s\n", block.Hash)
		fmt.Printf("Latest block miner: %s\n", block.Miner)
		fmt.Printf("Latest block transactions: %d\n", len(block.Transactions))
	}

	fmt.Println("\n=== Utility Functions Example ===")

	weiValue, _ := web3.ToWei("1", "ether")
	fmt.Printf("1 ETH in Wei: %s\n", weiValue.String())

	ethValue, _ := web3.FromWei(weiValue, "ether")
	fmt.Printf("Wei back to ETH: %s\n", ethValue)

	hexValue := web3.ToHex(12345)
	fmt.Printf("12345 in hex: %s\n", hexValue)

	isValidAddress := web3.IsAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	fmt.Printf("Is valid address: %t\n", isValidAddress)
}

func transactionExample() {
	client := web3.NewClient("https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY")

	fmt.Println("\n=== Go-Web3 Transaction Examples ===\n")

	// 1. Generate a new wallet
	fmt.Println("1. Creating a new wallet...")
	wallet, err := web3.CreateWallet(client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("   Address: %s\n", wallet.GetAddress())
	fmt.Printf("   Private Key: %s\n", wallet.GetPrivateKey())

	// 2. Transaction building and signing (without sending)
	fmt.Println("\n2. Building and signing a transaction...")
	txParams := web3.NewTransactionParams().
		SetTo("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045").
		SetValueInEther("0.1").
		SetGas(21000).
		SetGasPriceInGwei("20").
		SetNonce(0).
		SetChainID(big.NewInt(1)) // Mainnet

	privateKey, _ := web3.GeneratePrivateKey()
	signedTx, err := web3.SignTransaction(txParams, privateKey)
	if err != nil {
		log.Printf("   Error signing transaction: %v", err)
	} else {
		fmt.Printf("   Transaction Hash: %s\n", signedTx.Hash)
		fmt.Printf("   Raw Transaction: %s\n", signedTx.Raw[:66]+"...")
	}

	// 3. EIP-1559 Transaction
	fmt.Println("\n3. Building EIP-1559 transaction...")
	eip1559Params := web3.NewEIP1559TransactionParams()
	eip1559Params.To = "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
	eip1559Params.Value, _ = web3.ToWei("0.05", "ether")
	eip1559Params.Gas = 21000
	eip1559Params.MaxFeePerGas, _ = web3.ToWei("30", "gwei")
	eip1559Params.MaxPriorityFeePerGas, _ = web3.ToWei("2", "gwei")
	eip1559Params.Nonce = 1
	eip1559Params.ChainID = big.NewInt(1)

	signedEIP1559, err := web3.SignEIP1559Transaction(eip1559Params, privateKey)
	if err != nil {
		log.Printf("   Error signing EIP-1559 transaction: %v", err)
	} else {
		fmt.Printf("   EIP-1559 Hash: %s\n", signedEIP1559.Hash)
	}

	// 4. Contract interaction example
	fmt.Println("\n4. Contract interaction example...")

	// Example: ERC-20 transfer method call
	recipientAddress := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
	transferAmount := big.NewInt(1000000000000000000) // 1 token (18 decimals)

	// Encode transfer(address,uint256) function call
	transferData, err := web3.EncodeABI("transfer(address,uint256)", recipientAddress, transferAmount)
	if err != nil {
		log.Printf("   Error encoding ABI: %v", err)
	} else {
		fmt.Printf("   Transfer data: 0x%x\n", transferData[:20])
	}

	fmt.Println("\n=== Transaction examples completed! ===")
}

func main() {
	fmt.Println("=== Go-Web3 Library Examples ===")

	// Run basic blockchain query examples
	basicExample()

	// Run transaction examples
	transactionExample()
}
