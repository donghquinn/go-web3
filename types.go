package web3

import "math/big"

// Block parameter types
type BlockParameter string

const (
	BlockLatest   BlockParameter = "latest"
	BlockEarliest BlockParameter = "earliest"
	BlockPending  BlockParameter = "pending"
)

func (bp BlockParameter) String() string {
	return string(bp)
}

func BlockNumber(number uint64) BlockParameter {
	return BlockParameter(ToHex(number))
}

func BlockNumberBig(number *big.Int) BlockParameter {
	return BlockParameter(ToHex(number))
}

// Ether units enum
type EtherUnit string

const (
	Wei         EtherUnit = "wei"
	Kwei        EtherUnit = "kwei"
	Babbage     EtherUnit = "babbage"
	Femtoether  EtherUnit = "femtoether"
	Mwei        EtherUnit = "mwei"
	Lovelace    EtherUnit = "lovelace"
	Picoether   EtherUnit = "picoether"
	Gwei        EtherUnit = "gwei"
	Shannon     EtherUnit = "shannon"
	Nanoether   EtherUnit = "nanoether"
	Nano        EtherUnit = "nano"
	Szabo       EtherUnit = "szabo"
	Microether  EtherUnit = "microether"
	Micro       EtherUnit = "micro"
	Finney      EtherUnit = "finney"
	Milliether  EtherUnit = "milliether"
	Milli       EtherUnit = "milli"
	Ether       EtherUnit = "ether"
	EthUnit     EtherUnit = "eth"
	Kether      EtherUnit = "kether"
	Grand       EtherUnit = "grand"
	Mether      EtherUnit = "mether"
	Gether      EtherUnit = "gether"
	Tether      EtherUnit = "tether"
)

func (eu EtherUnit) String() string {
	return string(eu)
}

// Chain IDs for different networks
type ChainID uint64

const (
	ChainMainnet         ChainID = 1
	ChainGoerli          ChainID = 5
	ChainSepolia         ChainID = 11155111
	ChainOptimism        ChainID = 10
	ChainOptimismGoerli  ChainID = 420
	ChainArbitrum        ChainID = 42161
	ChainArbitrumGoerli  ChainID = 421613
	ChainPolygon         ChainID = 137
	ChainPolygonMumbai   ChainID = 80001
	ChainAvalanche       ChainID = 43114
	ChainAvalancheFuji   ChainID = 43113
	ChainBSC             ChainID = 56
	ChainBSCTestnet      ChainID = 97
	ChainFantom          ChainID = 250
	ChainFantomTestnet   ChainID = 4002
)

func (c ChainID) BigInt() *big.Int {
	return big.NewInt(int64(c))
}

func (c ChainID) Uint64() uint64 {
	return uint64(c)
}

// Transaction status
type TxStatus string

const (
	TxStatusSuccess TxStatus = "0x1"
	TxStatusFailure TxStatus = "0x0"
)

func (ts TxStatus) String() string {
	return string(ts)
}

func (ts TxStatus) IsSuccess() bool {
	return ts == TxStatusSuccess
}

func (ts TxStatus) IsFailure() bool {
	return ts == TxStatusFailure
}

// Gas price levels for optimization
type GasPriceLevel int

const (
	GasPriceSlow     GasPriceLevel = iota // Standard gas price
	GasPriceStandard                      // +10% buffer
	GasPriceFast                          // +20% buffer
	GasPriceRapid                         // +50% buffer
)

func (gpl GasPriceLevel) Multiplier() float64 {
	switch gpl {
	case GasPriceSlow:
		return 1.0
	case GasPriceStandard:
		return 1.1
	case GasPriceFast:
		return 1.2
	case GasPriceRapid:
		return 1.5
	default:
		return 1.0
	}
}

// Standard gas limits for common operations
type GasLimit uint64

const (
	GasLimitTransfer        GasLimit = 21000   // ETH transfer
	GasLimitTokenTransfer   GasLimit = 65000   // ERC-20 transfer
	GasLimitTokenApproval   GasLimit = 50000   // ERC-20 approval
	GasLimitContractCall    GasLimit = 100000  // Basic contract interaction
	GasLimitContractDeploy  GasLimit = 500000  // Contract deployment
	GasLimitComplexContract GasLimit = 1000000 // Complex contract operations
)

func (gl GasLimit) Uint64() uint64 {
	return uint64(gl)
}

// Common Ethereum addresses
type CommonAddress string

const (
	ZeroAddress     CommonAddress = "0x0000000000000000000000000000000000000000"
	BurnAddress     CommonAddress = "0x000000000000000000000000000000000000dEaD"
	NullAddress     CommonAddress = "0x0000000000000000000000000000000000000001"
	ENSRegistry     CommonAddress = "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"
	WETHMainnet     CommonAddress = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	USDCMainnet     CommonAddress = "0xA0b86a33E6417c48cd7a94Ca95e70aD2c51e74f7"
	USDTMainnet     CommonAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	DAIMainnet      CommonAddress = "0x6B175474E89094C44Da98b954EedeAC495271d0F"
	UniswapV3Router CommonAddress = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
)

func (ca CommonAddress) String() string {
	return string(ca)
}

// RPC method names
type RPCMethod string

const (
	EthGetBalance              RPCMethod = "eth_getBalance"
	EthGetBlockNumber          RPCMethod = "eth_blockNumber"
	EthGetGasPrice             RPCMethod = "eth_gasPrice"
	EthGetTransactionCount     RPCMethod = "eth_getTransactionCount"
	EthGetBlockByNumber        RPCMethod = "eth_getBlockByNumber"
	EthGetBlockByHash          RPCMethod = "eth_getBlockByHash"
	EthGetTransactionByHash    RPCMethod = "eth_getTransactionByHash"
	EthGetTransactionReceipt   RPCMethod = "eth_getTransactionReceipt"
	EthSendRawTransaction      RPCMethod = "eth_sendRawTransaction"
	EthEstimateGas             RPCMethod = "eth_estimateGas"
	EthCall                    RPCMethod = "eth_call"
	EthGetLogs                 RPCMethod = "eth_getLogs"
	EthGetStorageAt            RPCMethod = "eth_getStorageAt"
	EthGetCode                 RPCMethod = "eth_getCode"
	NetVersion                 RPCMethod = "net_version"
	Web3ClientVersion          RPCMethod = "web3_clientVersion"
	EthChainId                 RPCMethod = "eth_chainId"
	EthMaxPriorityFeePerGas    RPCMethod = "eth_maxPriorityFeePerGas"
	EthFeeHistory              RPCMethod = "eth_feeHistory"
)

func (rm RPCMethod) String() string {
	return string(rm)
}

// Common ERC-20 function signatures
type FunctionSignature string

const (
	FuncBalanceOf     FunctionSignature = "balanceOf(address)"
	FuncTransfer      FunctionSignature = "transfer(address,uint256)"
	FuncTransferFrom  FunctionSignature = "transferFrom(address,address,uint256)"
	FuncApprove       FunctionSignature = "approve(address,uint256)"
	FuncAllowance     FunctionSignature = "allowance(address,address)"
	FuncTotalSupply   FunctionSignature = "totalSupply()"
	FuncName          FunctionSignature = "name()"
	FuncSymbol        FunctionSignature = "symbol()"
	FuncDecimals      FunctionSignature = "decimals()"
	FuncOwner         FunctionSignature = "owner()"
	FuncMint          FunctionSignature = "mint(address,uint256)"
	FuncBurn          FunctionSignature = "burn(uint256)"
)

func (fs FunctionSignature) String() string {
	return string(fs)
}

// Network configurations
type NetworkConfig struct {
	Name     string
	ChainID  ChainID
	Currency string
	RPC      []string
	Explorer string
}

var Networks = map[ChainID]NetworkConfig{
	ChainMainnet: {
		Name:     "Ethereum Mainnet",
		ChainID:  ChainMainnet,
		Currency: "ETH",
		RPC: []string{
			"https://mainnet.infura.io/v3/YOUR_PROJECT_ID",
			"https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY",
			"https://rpc.ankr.com/eth",
		},
		Explorer: "https://etherscan.io",
	},
	ChainGoerli: {
		Name:     "Goerli Testnet",
		ChainID:  ChainGoerli,
		Currency: "GoerliETH",
		RPC: []string{
			"https://goerli.infura.io/v3/YOUR_PROJECT_ID",
			"https://eth-goerli.alchemyapi.io/v2/YOUR_API_KEY",
		},
		Explorer: "https://goerli.etherscan.io",
	},
	ChainSepolia: {
		Name:     "Sepolia Testnet",
		ChainID:  ChainSepolia,
		Currency: "SepoliaETH",
		RPC: []string{
			"https://sepolia.infura.io/v3/YOUR_PROJECT_ID",
			"https://eth-sepolia.alchemyapi.io/v2/YOUR_API_KEY",
		},
		Explorer: "https://sepolia.etherscan.io",
	},
	ChainPolygon: {
		Name:     "Polygon",
		ChainID:  ChainPolygon,
		Currency: "MATIC",
		RPC: []string{
			"https://polygon-mainnet.infura.io/v3/YOUR_PROJECT_ID",
			"https://polygon-mainnet.g.alchemy.com/v2/YOUR_API_KEY",
		},
		Explorer: "https://polygonscan.com",
	},
}