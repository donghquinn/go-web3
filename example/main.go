package main

import (
	"context"
	"fmt"
	"log"

	"github.com/donghquinn/go-web3"
)

func main() {
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