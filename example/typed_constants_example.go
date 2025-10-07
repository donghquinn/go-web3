package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/donghquinn/go-web3"
)

func typedConstantsExample() {
	client := web3.NewClient("https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY")
	ctx := context.Background()

	fmt.Printf("\n=== Typed Constants Example ===\n")

	// 1. Using typed block parameters
	fmt.Println("1. Block Parameters:")
	address := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"

	// Different ways to specify blocks
	latestBalance, _ := client.Eth().GetBalance(ctx, address, web3.BlockLatest)
	fmt.Printf("   Latest balance: %s\n", latestBalance.String())

	pendingNonce, _ := client.Eth().GetTransactionCount(ctx, address, web3.BlockPending)
	fmt.Printf("   Pending nonce: %d\n", pendingNonce)

	// Using specific block number
	blockNum := web3.BlockNumber(18500000)
	historicalBalance, err := client.Eth().GetBalance(ctx, address, blockNum)
	if err == nil {
		fmt.Printf("   Block #18500000 balance: %s\n", historicalBalance.String())
	}

	// 2. Using typed ether units
	fmt.Println("\n2. Ether Unit Conversions:")

	// Convert using typed units
	oneEth, _ := web3.ToWei("1", web3.Ether)
	oneGwei, _ := web3.ToWei("1", web3.Gwei)
	oneSzabo, _ := web3.ToWei("1", web3.Szabo)

	fmt.Printf("   1 ETH = %s Wei\n", oneEth.String())
	fmt.Printf("   1 Gwei = %s Wei\n", oneGwei.String())
	fmt.Printf("   1 Szabo = %s Wei\n", oneSzabo.String())

	// Using helper functions
	ethValue, _ := web3.EtherToWei("2.5")
	gweiValue, _ := web3.GweiToWei("20")

	fmt.Printf("   2.5 ETH = %s Wei\n", ethValue.String())
	fmt.Printf("   20 Gwei = %s Wei\n", gweiValue.String())

	// Convert back
	ethDisplay, _ := web3.WeiToEther(ethValue)
	gweiDisplay, _ := web3.WeiToGwei(gweiValue)

	fmt.Printf("   %s Wei = %s ETH\n", ethValue.String(), ethDisplay)
	fmt.Printf("   %s Wei = %s Gwei\n", gweiValue.String(), gweiDisplay)

	// 3. Using typed chain IDs
	fmt.Println("\n3. Chain ID Examples:")

	chains := []web3.ChainID{
		web3.ChainMainnet,
		web3.ChainGoerli,
		web3.ChainPolygon,
		web3.ChainArbitrum,
	}

	for _, chainID := range chains {
		config, err := web3.GetNetworkConfig(chainID)
		if err == nil {
			fmt.Printf("   Chain %d: %s (%s)\n", chainID.Uint64(), config.Name, config.Currency)
			fmt.Printf("     Testnet: %t\n", web3.IsTestnet(chainID))
		}
	}

	// 4. Using typed gas limits
	fmt.Println("\n4. Gas Limit Constants:")

	fmt.Printf("   ETH Transfer: %d gas\n", web3.GasLimitTransfer.Uint64())
	fmt.Printf("   Token Transfer: %d gas\n", web3.GasLimitTokenTransfer.Uint64())
	fmt.Printf("   Token Approval: %d gas\n", web3.GasLimitTokenApproval.Uint64())
	fmt.Printf("   Contract Deployment: %d gas\n", web3.GasLimitContractDeploy.Uint64())

	// 5. Using gas price levels
	fmt.Println("\n5. Gas Price Optimization:")

	_, err = client.Eth().GetGasPrice(ctx)
	if err == nil {
		levels := []web3.GasPriceLevel{
			web3.GasPriceSlow,
			web3.GasPriceStandard,
			web3.GasPriceFast,
			web3.GasPriceRapid,
		}

		for _, level := range levels {
			optimal, err := web3.GetOptimalGasPrice(ctx, client, level)
			if err == nil {
				gweiPrice, _ := web3.WeiToGwei(optimal)
				fmt.Printf("   %s: %s Gwei (%.1fx multiplier)\n",
					getLevelName(level), gweiPrice, level.Multiplier())
			}
		}
	}

	// 6. Transaction helpers
	fmt.Println("\n6. Transaction Builder Helpers:")

	// Simple ETH transfer
	ethTransfer := web3.NewSimpleTransfer(
		"0xRecipientAddress",
		"0.1",
		web3.ChainMainnet,
	)
	fmt.Printf("   ETH Transfer - Gas: %d, Chain: %d\n",
		ethTransfer.Gas, ethTransfer.ChainID.Uint64())

	// Token transfer
	tokenAmount := big.NewInt(1000000000000000000) // 1 token with 18 decimals
	tokenTransfer, err := web3.NewTokenTransfer(
		web3.USDCMainnet.String(),
		"0xRecipientAddress",
		tokenAmount,
		web3.ChainMainnet,
	)
	if err == nil {
		fmt.Printf("   Token Transfer - Gas: %d, To: %s\n",
			tokenTransfer.Gas, tokenTransfer.To)
	}

	// 7. Common addresses
	fmt.Println("\n7. Common Addresses:")

	fmt.Printf("   Zero Address: %s\n", web3.ZeroAddress.String())
	fmt.Printf("   Burn Address: %s\n", web3.BurnAddress.String())
	fmt.Printf("   WETH (Mainnet): %s\n", web3.WETHMainnet.String())
	fmt.Printf("   USDC (Mainnet): %s\n", web3.USDCMainnet.String())

	// Address checks
	fmt.Printf("   Is zero address: %t\n", web3.IsZeroAddress(web3.ZeroAddress.String()))
	fmt.Printf("   Is burn address: %t\n", web3.IsBurnAddress(web3.BurnAddress.String()))

	// 8. Function signatures
	fmt.Println("\n8. Common ERC-20 Function Signatures:")

	signatures := []web3.FunctionSignature{
		web3.FuncBalanceOf,
		web3.FuncTransfer,
		web3.FuncApprove,
		web3.FuncAllowance,
	}

	for _, sig := range signatures {
		fmt.Printf("   %s\n", sig.String())
	}

	fmt.Println("\n=== Typed constants example completed! ===")
}

func getLevelName(level web3.GasPriceLevel) string {
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
