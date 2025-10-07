package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/donghquinn/go-web3"
)

func blockchainHelperExample() {
	client := web3.NewClient("https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY")
	ctx := context.Background()

	fmt.Printf("\n=== Go-Blockchain-Helper Integration Examples ===\n")

	// 1. Enhanced Unit Conversions
	fmt.Println("1. Enhanced Unit Conversions with go-blockchain-helper:")

	// Using the new ParseEther function
	ethAmount, err := web3.ParseEther("2.5")
	if err != nil {
		log.Printf("Error parsing ether: %v", err)
	} else {
		fmt.Printf("   Parsed 2.5 ETH: %s Wei\n", ethAmount.String())
	}

	// Using FormatEther with custom decimals
	weiAmount := big.NewInt(1500000000000000000) // 1.5 ETH
	formattedEth := web3.FormatEther(weiAmount, 2)
	fmt.Printf("   Formatted Wei to ETH (2 decimals): %s\n", formattedEth)

	// Custom unit parsing
	tokenAmount, err := web3.ParseUnits("100.5", 6) // USDC has 6 decimals
	if err != nil {
		log.Printf("Error parsing units: %v", err)
	} else {
		fmt.Printf("   Parsed 100.5 USDC (6 decimals): %s\n", tokenAmount.String())
	}

	// 2. Enhanced ERC20 Token Operations
	fmt.Println("\n2. Enhanced ERC20 Token Operations:")

	// Create ERC20 token instance with go-blockchain-helper
	usdcToken := web3.NewERC20Token(web3.USDCMainnet.String(), "USD Coin", "USDC", 6)

	// Enhanced token transfer encoding
	recipient := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
	transferAmount := big.NewInt(1000000) // 1 USDC (6 decimals)

	transferData, err := web3.EncodeERC20Transfer(usdcToken, recipient, transferAmount)
	if err != nil {
		log.Printf("Error encoding transfer: %v", err)
	} else {
		fmt.Printf("   Transfer data: 0x%x...\n", transferData[:20])
	}

	// Enhanced token approval encoding
	spender := "0x1234567890123456789012345678901234567890"
	approveAmount := big.NewInt(10000000) // 10 USDC

	approveData, err := web3.EncodeERC20Approve(usdcToken, spender, approveAmount)
	if err != nil {
		log.Printf("Error encoding approval: %v", err)
	} else {
		fmt.Printf("   Approval data: 0x%x...\n", approveData[:20])
	}

	// 3. Enhanced Transaction Creation
	fmt.Println("\n3. Enhanced Transaction Creation:")

	// Create transaction with enhanced estimation
	ethValue, _ := web3.ToWei("0.1", web3.Ether)
	txParams, err := web3.CreateTransactionWithEstimate(
		recipient,
		ethValue,
		[]byte{},
		web3.ChainMainnet,
	)
	if err != nil {
		log.Printf("Error creating transaction: %v", err)
	} else {
		fmt.Printf("   Enhanced transaction created with gas: %d\n", txParams.Gas)
		gasPriceGwei, _ := web3.WeiToGwei(txParams.GasPrice)
		fmt.Printf("   Gas price: %s Gwei\n", gasPriceGwei)
	}

	// 4. Enhanced Token Contract Interactions
	fmt.Println("\n4. Enhanced Token Contract Interactions:")

	// Get token balance using go-blockchain-helper
	address := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
	balance, err := web3.GetTokenBalance(ctx, client, web3.USDCMainnet.String(), address)
	if err != nil {
		log.Printf("Error getting token balance: %v", err)
	} else {
		// Format balance with USDC decimals
		balanceFormatted := web3.FormatUnits(balance, 6)
		fmt.Printf("   USDC Balance: %s USDC\n", balanceFormatted)
	}

	// Get token allowance using go-blockchain-helper
	allowance, err := web3.GetTokenAllowance(ctx, client, web3.USDCMainnet.String(), address, spender)
	if err != nil {
		log.Printf("Error getting allowance: %v", err)
	} else {
		allowanceFormatted := web3.FormatUnits(allowance, 6)
		fmt.Printf("   USDC Allowance: %s USDC\n", allowanceFormatted)
	}

	// 5. ERC721 (NFT) Operations
	fmt.Println("\n5. Enhanced ERC721 (NFT) Operations:")

	// Create ERC721 token instance
	nftContract := "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D" // Bored Ape Yacht Club
	nftToken := web3.NewERC721Token(nftContract, "BoredApeYachtClub", "BAYC")

	// Encode NFT transfer
	from := "0x1234567890123456789012345678901234567890"
	to := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
	tokenId := big.NewInt(1234)

	nftTransferData, err := web3.EncodeERC721Transfer(nftToken, from, to, tokenId)
	if err != nil {
		log.Printf("Error encoding NFT transfer: %v", err)
	} else {
		fmt.Printf("   NFT Transfer data: 0x%x...\n", nftTransferData[:20])
	}

	// 6. Enhanced Gas Price Optimization
	fmt.Println("\n6. Enhanced Gas Price Optimization:")

	// Get optimal gas prices using different levels
	levels := []web3.GasPriceLevel{
		web3.GasPriceSlow,
		web3.GasPriceStandard,
		web3.GasPriceFast,
		web3.GasPriceRapid,
	}

	for _, level := range levels {
		optimalGas, err := web3.GetOptimalGasPrice(ctx, client, level)
		if err != nil {
			log.Printf("Error getting optimal gas: %v", err)
			continue
		}

		gweiPrice, _ := web3.WeiToGwei(optimalGas)
		levelName := getGasPriceLevelName(level)
		fmt.Printf("   %s gas price: %s Gwei (%.1fx)\n",
			levelName, gweiPrice, level.Multiplier())
	}

	// 7. Transaction Fee Calculation
	fmt.Println("\n7. Enhanced Transaction Fee Calculation:")

	gasLimit := uint64(21000)
	gasPrice := big.NewInt(20000000000) // 20 Gwei

	fee := web3.CalculateTransactionFee(gasLimit, gasPrice)
	feeEth, _ := web3.WeiToEther(fee)
	fmt.Printf("   Transaction fee: %s ETH\n", feeEth)

	// 8. Event Processing Capabilities
	fmt.Println("\n8. Event Processing with go-blockchain-helper:")

	// Create event monitor
	eventMonitor := web3.CreateEventMonitor()
	fmt.Printf("   Event monitor created: %v\n", eventMonitor != nil)

	// 9. Address Validation Enhancement
	fmt.Println("\n9. Enhanced Address Validation:")

	testAddresses := []string{
		"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
		"0xinvalid",
		web3.ZeroAddress.String(),
		web3.BurnAddress.String(),
	}

	for _, addr := range testAddresses {
		isValid := web3.ValidateAddress(addr)
		isZero := web3.IsZeroAddress(addr)
		isBurn := web3.IsBurnAddress(addr)

		fmt.Printf("   %s - Valid: %t, Zero: %t, Burn: %t\n",
			addr, isValid, isZero, isBurn)
	}

	// 10. Simplified Transaction Building
	fmt.Println("\n10. Simplified Transaction Building:")

	// Simple ETH transfer
	simpleTransfer := web3.NewSimpleTransfer(
		"0xRecipient",
		"1.0",
		web3.ChainMainnet,
	)
	fmt.Printf("   Simple transfer gas limit: %d\n", simpleTransfer.Gas)

	// Token transfer transaction
	tokenTransferTx, err := web3.NewTokenTransfer(
		web3.USDCMainnet.String(),
		"0xRecipient",
		big.NewInt(1000000), // 1 USDC
		web3.ChainMainnet,
	)
	if err != nil {
		log.Printf("Error creating token transfer: %v", err)
	} else {
		fmt.Printf("   Token transfer gas limit: %d\n", tokenTransferTx.Gas)
		fmt.Printf("   Token transfer data length: %d bytes\n", len(tokenTransferTx.Data))
	}

	fmt.Println("\n=== Go-blockchain-helper integration example completed! ===")
}

func getGasPriceLevelName(level web3.GasPriceLevel) string {
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
