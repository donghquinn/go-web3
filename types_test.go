package web3

import (
	"math/big"
	"testing"
)

func TestBlockParameter(t *testing.T) {
	if BlockLatest.String() != "latest" {
		t.Errorf("BlockLatest = %q, want %q", BlockLatest.String(), "latest")
	}
	if BlockEarliest.String() != "earliest" {
		t.Errorf("BlockEarliest = %q, want %q", BlockEarliest.String(), "earliest")
	}
	if BlockPending.String() != "pending" {
		t.Errorf("BlockPending = %q, want %q", BlockPending.String(), "pending")
	}
}

func TestBlockNumber(t *testing.T) {
	bp := BlockNumber(100)
	if bp.String() != "0x64" {
		t.Errorf("BlockNumber(100) = %q, want %q", bp.String(), "0x64")
	}
}

func TestBlockNumberBig(t *testing.T) {
	bp := BlockNumberBig(big.NewInt(255))
	if bp.String() != "0xff" {
		t.Errorf("BlockNumberBig(255) = %q, want %q", bp.String(), "0xff")
	}
}

func TestEtherUnitString(t *testing.T) {
	units := []struct {
		unit EtherUnit
		want string
	}{
		{Wei, "wei"},
		{Gwei, "gwei"},
		{Ether, "ether"},
		{EthUnit, "eth"},
		{Kether, "kether"},
		{Tether, "tether"},
	}
	for _, tt := range units {
		if tt.unit.String() != tt.want {
			t.Errorf("EtherUnit %q String() = %q, want %q", tt.unit, tt.unit.String(), tt.want)
		}
	}
}

func TestChainIDBigInt(t *testing.T) {
	if ChainMainnet.BigInt().Int64() != 1 {
		t.Errorf("ChainMainnet.BigInt() = %v, want 1", ChainMainnet.BigInt())
	}
	if ChainPolygon.BigInt().Int64() != 137 {
		t.Errorf("ChainPolygon.BigInt() = %v, want 137", ChainPolygon.BigInt())
	}
}

func TestChainIDUint64(t *testing.T) {
	if ChainMainnet.Uint64() != 1 {
		t.Errorf("ChainMainnet.Uint64() = %v, want 1", ChainMainnet.Uint64())
	}
	if ChainSepolia.Uint64() != 11155111 {
		t.Errorf("ChainSepolia.Uint64() = %v, want 11155111", ChainSepolia.Uint64())
	}
}

func TestTxStatus(t *testing.T) {
	if !TxStatusSuccess.IsSuccess() {
		t.Error("TxStatusSuccess.IsSuccess() should be true")
	}
	if TxStatusSuccess.IsFailure() {
		t.Error("TxStatusSuccess.IsFailure() should be false")
	}
	if !TxStatusFailure.IsFailure() {
		t.Error("TxStatusFailure.IsFailure() should be true")
	}
	if TxStatusFailure.IsSuccess() {
		t.Error("TxStatusFailure.IsSuccess() should be false")
	}
	if TxStatusSuccess.String() != "0x1" {
		t.Errorf("TxStatusSuccess.String() = %q, want %q", TxStatusSuccess.String(), "0x1")
	}
	if TxStatusFailure.String() != "0x0" {
		t.Errorf("TxStatusFailure.String() = %q, want %q", TxStatusFailure.String(), "0x0")
	}
}

func TestGasPriceLevelMultiplier(t *testing.T) {
	tests := []struct {
		level GasPriceLevel
		want  float64
	}{
		{GasPriceSlow, 1.0},
		{GasPriceStandard, 1.1},
		{GasPriceFast, 1.2},
		{GasPriceRapid, 1.5},
		{GasPriceLevel(99), 1.0}, // unknown defaults to 1.0
	}
	for _, tt := range tests {
		got := tt.level.Multiplier()
		if got != tt.want {
			t.Errorf("GasPriceLevel(%d).Multiplier() = %v, want %v", tt.level, got, tt.want)
		}
	}
}

func TestGasLimitValues(t *testing.T) {
	if GasLimitTransfer.Uint64() != 21000 {
		t.Errorf("GasLimitTransfer = %d, want 21000", GasLimitTransfer.Uint64())
	}
	if GasLimitTokenTransfer.Uint64() != 65000 {
		t.Errorf("GasLimitTokenTransfer = %d, want 65000", GasLimitTokenTransfer.Uint64())
	}
	if GasLimitTokenApproval.Uint64() != 50000 {
		t.Errorf("GasLimitTokenApproval = %d, want 50000", GasLimitTokenApproval.Uint64())
	}
	if GasLimitContractDeploy.Uint64() != 500000 {
		t.Errorf("GasLimitContractDeploy = %d, want 500000", GasLimitContractDeploy.Uint64())
	}
}

func TestCommonAddressString(t *testing.T) {
	if ZeroAddress.String() != "0x0000000000000000000000000000000000000000" {
		t.Errorf("ZeroAddress.String() unexpected: %q", ZeroAddress.String())
	}
	if BurnAddress.String() != "0x000000000000000000000000000000000000dEaD" {
		t.Errorf("BurnAddress.String() unexpected: %q", BurnAddress.String())
	}
}

func TestRPCMethodString(t *testing.T) {
	tests := []struct {
		method RPCMethod
		want   string
	}{
		{EthGetBalance, "eth_getBalance"},
		{EthGetBlockNumber, "eth_blockNumber"},
		{EthSendRawTransaction, "eth_sendRawTransaction"},
		{EthChainId, "eth_chainId"},
		{NetVersion, "net_version"},
	}
	for _, tt := range tests {
		if tt.method.String() != tt.want {
			t.Errorf("RPCMethod.String() = %q, want %q", tt.method.String(), tt.want)
		}
	}
}

func TestFunctionSignatureString(t *testing.T) {
	if FuncTransfer.String() != "transfer(address,uint256)" {
		t.Errorf("FuncTransfer.String() = %q", FuncTransfer.String())
	}
	if FuncBalanceOf.String() != "balanceOf(address)" {
		t.Errorf("FuncBalanceOf.String() = %q", FuncBalanceOf.String())
	}
}

func TestNetworksMap(t *testing.T) {
	cfg, ok := Networks[ChainMainnet]
	if !ok {
		t.Fatal("ChainMainnet not in Networks")
	}
	if cfg.Name != "Ethereum Mainnet" {
		t.Errorf("Networks[ChainMainnet].Name = %q, want %q", cfg.Name, "Ethereum Mainnet")
	}
	if cfg.Currency != "ETH" {
		t.Errorf("Networks[ChainMainnet].Currency = %q, want %q", cfg.Currency, "ETH")
	}
	if len(cfg.RPC) == 0 {
		t.Error("Networks[ChainMainnet].RPC is empty")
	}
}
