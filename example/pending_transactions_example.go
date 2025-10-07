package main

import (
	"context"
	"fmt"
	"log"

	"github.com/donghquinn/go-web3"
)

func pendingTransactionsExample() {
	client := web3.NewClient("https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY")
	ctx := context.Background()

	fmt.Printf("\n=== Pending Transactions Example ===\n")

	// 1. Get all pending transactions
	fmt.Println("1. Getting all pending transactions...")
	pendingTxs, err := client.Eth().GetPendingTransactions(ctx)
	if err != nil {
		log.Printf("Error getting pending transactions: %v", err)
		return
	}

	fmt.Printf("   Found %d pending transactions\n", len(pendingTxs))
	
	// Display first few transactions
	displayCount := 3
	if len(pendingTxs) < displayCount {
		displayCount = len(pendingTxs)
	}
	
	for i := 0; i < displayCount; i++ {
		tx := pendingTxs[i]
		fmt.Printf("   Transaction %d:\n", i+1)
		fmt.Printf("     Hash: %s\n", tx.Hash)
		fmt.Printf("     From: %s\n", tx.From)
		fmt.Printf("     To: %s\n", tx.To)
		fmt.Printf("     Value: %s Wei\n", tx.Value)
		fmt.Printf("     Gas Price: %s Wei\n", tx.GasPrice)
		fmt.Printf("     Gas Limit: %s\n", tx.Gas)
		fmt.Printf("     Nonce: %s\n", tx.Nonce)
		
		// Convert gas price to Gwei for readability
		if tx.GasPrice != "" && tx.GasPrice != "0x0" {
			gasPriceWei, parseErr := web3.FromHex(tx.GasPrice)
			if parseErr == nil {
				gasPriceGwei, _ := web3.WeiToGwei(gasPriceWei)
				fmt.Printf("     Gas Price: %s Gwei\n", gasPriceGwei)
			}
		}
		fmt.Println()
	}

	// 2. Get pending transaction count
	fmt.Println("2. Getting pending transaction count...")
	pendingCount, err := client.Eth().GetPendingTransactionCount(ctx)
	if err != nil {
		log.Printf("Error getting pending transaction count: %v", err)
	} else {
		fmt.Printf("   Total pending transactions: %d\n", pendingCount)
	}

	// 3. Get pending transactions for specific account
	fmt.Println("\n3. Getting pending transactions for specific account...")
	
	// Use a popular address (Uniswap V3 Router)
	targetAddress := "0xE592427A0AEce92De3Edee1F18E0157C05861564"
	accountTxs, err := client.Eth().GetAccountPendingTransactions(ctx, targetAddress)
	if err != nil {
		log.Printf("Error getting account pending transactions: %v", err)
	} else {
		fmt.Printf("   Found %d pending transactions for address %s\n", 
			len(accountTxs), targetAddress)
		
		for i, tx := range accountTxs {
			fmt.Printf("   Transaction %d:\n", i+1)
			fmt.Printf("     Hash: %s\n", tx.Hash)
			fmt.Printf("     From: %s\n", tx.From)
			fmt.Printf("     To: %s\n", tx.To)
			fmt.Printf("     Value: %s\n", tx.Value)
		}
	}

	// 4. Check if specific transaction is pending
	fmt.Println("\n4. Checking if specific transaction is pending...")
	
	if len(pendingTxs) > 0 {
		testTxHash := pendingTxs[0].Hash
		isPending, err := client.Eth().IsPendingTransaction(ctx, testTxHash)
		if err != nil {
			log.Printf("Error checking if transaction is pending: %v", err)
		} else {
			fmt.Printf("   Transaction %s is pending: %t\n", testTxHash, isPending)
		}
	} else {
		// Test with a known non-pending transaction hash
		testTxHash := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		isPending, err := client.Eth().IsPendingTransaction(ctx, testTxHash)
		if err != nil {
			log.Printf("Error checking if transaction is pending: %v", err)
		} else {
			fmt.Printf("   Transaction %s is pending: %t\n", testTxHash, isPending)
		}
	}

	// 5. Monitor pending transactions with gas price analysis
	fmt.Println("\n5. Analyzing gas prices in pending transactions...")
	
	if len(pendingTxs) > 0 {
		var totalGasPrice float64
		var validTxCount int
		
		for _, tx := range pendingTxs {
			if tx.GasPrice != "" && tx.GasPrice != "0x0" {
				gasPriceWei, parseErr := web3.FromHex(tx.GasPrice)
				if parseErr == nil {
					gasPriceGwei, _ := web3.WeiToGwei(gasPriceWei)
					if gasPriceFloat, parseErr2 := parseFloat(gasPriceGwei); parseErr2 == nil {
						totalGasPrice += gasPriceFloat
						validTxCount++
					}
				}
			}
		}
		
		if validTxCount > 0 {
			avgGasPrice := totalGasPrice / float64(validTxCount)
			fmt.Printf("   Average gas price in pending pool: %.2f Gwei\n", avgGasPrice)
			fmt.Printf("   Based on %d transactions with valid gas prices\n", validTxCount)
		}
	}

	// 6. Compare with current network gas price
	fmt.Println("\n6. Comparing with current network gas price...")
	
	currentGasPrice, err := client.Eth().GetGasPrice(ctx)
	if err != nil {
		log.Printf("Error getting current gas price: %v", err)
	} else {
		currentGwei, _ := web3.WeiToGwei(currentGasPrice)
		fmt.Printf("   Current network gas price: %s Gwei\n", currentGwei)
		
		// Get optimal gas prices for different levels
		levels := []web3.GasPriceLevel{
			web3.GasPriceSlow,
			web3.GasPriceStandard,
			web3.GasPriceFast,
			web3.GasPriceRapid,
		}
		
		fmt.Println("   Recommended gas prices:")
		for _, level := range levels {
			optimal, err := web3.GetOptimalGasPrice(ctx, client, level)
			if err == nil {
				optimalGwei, _ := web3.WeiToGwei(optimal)
				fmt.Printf("     %s: %s Gwei (%.1fx)\n", 
					getPendingGasPriceLevelName(level), optimalGwei, level.Multiplier())
			}
		}
	}

	fmt.Println("\n=== Pending transactions example completed! ===")
}

func getPendingGasPriceLevelName(level web3.GasPriceLevel) string {
	switch level {
	case web3.GasPriceSlow:
		return "Slow    "
	case web3.GasPriceStandard:
		return "Standard"
	case web3.GasPriceFast:
		return "Fast    "
	case web3.GasPriceRapid:
		return "Rapid   "
	default:
		return "Unknown "
	}
}

// Helper function to parse float from string
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}