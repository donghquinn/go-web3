package web3

import (
	"context"
	"math/big"
	"testing"
)

func TestEtherToWei(t *testing.T) {
	wei, err := EtherToWei("1")
	if err != nil {
		t.Fatalf("EtherToWei() error = %v", err)
	}
	expected, _ := new(big.Int).SetString("1000000000000000000", 10)
	if wei.Cmp(expected) != 0 {
		t.Errorf("EtherToWei(1) = %v, want %v", wei, expected)
	}
}

func TestWeiToEther(t *testing.T) {
	wei, _ := new(big.Int).SetString("1000000000000000000", 10)
	result, err := WeiToEther(wei)
	if err != nil {
		t.Fatalf("WeiToEther() error = %v", err)
	}
	if result == "" {
		t.Error("WeiToEther() returned empty string")
	}
}

func TestGweiToWei(t *testing.T) {
	wei, err := GweiToWei("1")
	if err != nil {
		t.Fatalf("GweiToWei() error = %v", err)
	}
	expected := big.NewInt(1000000000)
	if wei.Cmp(expected) != 0 {
		t.Errorf("GweiToWei(1) = %v, want %v", wei, expected)
	}
}

func TestWeiToGwei(t *testing.T) {
	wei := big.NewInt(1000000000) // 1 gwei
	result, err := WeiToGwei(wei)
	if err != nil {
		t.Fatalf("WeiToGwei() error = %v", err)
	}
	if result == "" {
		t.Error("WeiToGwei() returned empty string")
	}
}

func TestParseEther(t *testing.T) {
	wei, err := ParseEther("1")
	if err != nil {
		t.Fatalf("ParseEther() error = %v", err)
	}
	expected, _ := new(big.Int).SetString("1000000000000000000", 10)
	if wei.Cmp(expected) != 0 {
		t.Errorf("ParseEther(1) = %v, want %v", wei, expected)
	}
}

func TestFormatEther(t *testing.T) {
	wei, _ := new(big.Int).SetString("1000000000000000000", 10)
	result := FormatEther(wei, 18)
	if result == "" {
		t.Error("FormatEther() returned empty string")
	}
}

func TestParseUnits(t *testing.T) {
	val, err := ParseUnits("1", 6) // USDC-style 6 decimals
	if err != nil {
		t.Fatalf("ParseUnits() error = %v", err)
	}
	expected := big.NewInt(1000000)
	if val.Cmp(expected) != 0 {
		t.Errorf("ParseUnits(1, 6) = %v, want %v", val, expected)
	}
}

func TestFormatUnits(t *testing.T) {
	val := big.NewInt(1000000)
	result := FormatUnits(val, 6)
	if result == "" {
		t.Error("FormatUnits() returned empty string")
	}
}

func TestGetNetworkConfig(t *testing.T) {
	cfg, err := GetNetworkConfig(ChainMainnet)
	if err != nil {
		t.Fatalf("GetNetworkConfig(ChainMainnet) error = %v", err)
	}
	if cfg.Name != "Ethereum Mainnet" {
		t.Errorf("Name = %q, want %q", cfg.Name, "Ethereum Mainnet")
	}
}

func TestGetNetworkConfigUnknown(t *testing.T) {
	_, err := GetNetworkConfig(ChainID(999999))
	if err == nil {
		t.Error("expected error for unknown chain ID")
	}
}

func TestIsTestnet(t *testing.T) {
	testnets := []ChainID{
		ChainGoerli, ChainSepolia, ChainOptimismGoerli,
		ChainArbitrumGoerli, ChainPolygonMumbai,
		ChainAvalancheFuji, ChainBSCTestnet, ChainFantomTestnet,
	}
	for _, id := range testnets {
		if !IsTestnet(id) {
			t.Errorf("IsTestnet(%d) = false, want true", id)
		}
	}
	if IsTestnet(ChainMainnet) {
		t.Error("IsTestnet(ChainMainnet) = true, want false")
	}
}

func TestIsMainnet(t *testing.T) {
	if !IsMainnet(ChainMainnet) {
		t.Error("IsMainnet(ChainMainnet) = false, want true")
	}
	if IsMainnet(ChainGoerli) {
		t.Error("IsMainnet(ChainGoerli) = true, want false")
	}
}

func TestNewSimpleTransfer(t *testing.T) {
	to := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	tp := NewSimpleTransfer(to, "1", ChainMainnet)
	if tp == nil {
		t.Fatal("NewSimpleTransfer() returned nil")
	}
	if tp.To != to {
		t.Errorf("To = %q, want %q", tp.To, to)
	}
	if tp.Gas != GasLimitTransfer.Uint64() {
		t.Errorf("Gas = %d, want %d", tp.Gas, GasLimitTransfer.Uint64())
	}
	expected, _ := EtherToWei("1")
	if tp.Value.Cmp(expected) != 0 {
		t.Errorf("Value = %v, want %v", tp.Value, expected)
	}
}

func TestIsZeroAddress(t *testing.T) {
	if !IsZeroAddress(ZeroAddress.String()) {
		t.Error("IsZeroAddress(ZeroAddress) = false, want true")
	}
	if !IsZeroAddress("0x0") {
		t.Error("IsZeroAddress(0x0) = false, want true")
	}
	if IsZeroAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e") {
		t.Error("IsZeroAddress(non-zero) = true, want false")
	}
}

func TestIsBurnAddress(t *testing.T) {
	if !IsBurnAddress(BurnAddress.String()) {
		t.Error("IsBurnAddress(BurnAddress) = false, want true")
	}
	if IsBurnAddress(ZeroAddress.String()) {
		t.Error("IsBurnAddress(ZeroAddress) = true, want false")
	}
}

func TestCalculateTransactionFee(t *testing.T) {
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(20000000000) // 20 gwei
	fee := CalculateTransactionFee(gasLimit, gasPrice)

	expected := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)
	if fee.Cmp(expected) != 0 {
		t.Errorf("CalculateTransactionFee() = %v, want %v", fee, expected)
	}
}

func TestIsTransactionSuccess(t *testing.T) {
	receipt := &TransactionReceipt{Status: "0x1"}
	if !IsTransactionSuccess(receipt) {
		t.Error("IsTransactionSuccess(0x1) = false, want true")
	}
	if IsTransactionFailure(receipt) {
		t.Error("IsTransactionFailure(0x1) = true, want false")
	}
}

func TestIsTransactionFailure(t *testing.T) {
	receipt := &TransactionReceipt{Status: "0x0"}
	if !IsTransactionFailure(receipt) {
		t.Error("IsTransactionFailure(0x0) = false, want true")
	}
	if IsTransactionSuccess(receipt) {
		t.Error("IsTransactionSuccess(0x0) = true, want false")
	}
}

func TestGetOptimalGasPrice(t *testing.T) {
	srv := newMockRPCServer("0x4a817c800") // 20 gwei
	defer srv.Close()

	client := NewClient(srv.URL)
	price, err := GetOptimalGasPrice(context.Background(), client, GasPriceFast)
	if err != nil {
		t.Fatalf("GetOptimalGasPrice() error = %v", err)
	}
	// 20 gwei * 1.2 ≈ 24 gwei (big.Float may truncate by 1)
	base, _ := ToWei("20", Gwei)
	if price.Cmp(base) <= 0 {
		t.Errorf("GetOptimalGasPrice(Fast) = %v should be greater than base %v", price, base)
	}
	rapid, _ := GetOptimalGasPrice(context.Background(), client, GasPriceRapid)
	if rapid.Cmp(price) <= 0 {
		t.Errorf("Rapid price %v should be greater than Fast price %v", rapid, price)
	}
}

func TestEstimateGasWithBuffer(t *testing.T) {
	srv := newMockRPCServer("0x5208") // 21000
	defer srv.Close()

	client := NewClient(srv.URL)
	gas, err := EstimateGasWithBuffer(context.Background(), client, map[string]interface{}{
		"to": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
	}, 0.1)
	if err != nil {
		t.Fatalf("EstimateGasWithBuffer() error = %v", err)
	}
	// 21000 + 21000*0.1 = 23100
	if gas != 23100 {
		t.Errorf("EstimateGasWithBuffer() = %d, want 23100", gas)
	}
}

func TestNewERC20Token(t *testing.T) {
	token := NewERC20Token("0xcontract", "TestToken", "TST", 18)
	if token == nil {
		t.Fatal("NewERC20Token() returned nil")
	}
}

func TestEncodeERC20Transfer(t *testing.T) {
	token := NewERC20Token("0xcontract", "Token", "TKN", 18)
	data, err := EncodeERC20Transfer(token, "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", big.NewInt(1000))
	if err != nil {
		t.Fatalf("EncodeERC20Transfer() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeERC20Transfer() returned empty data")
	}
}

func TestEncodeERC20Approve(t *testing.T) {
	token := NewERC20Token("0xcontract", "Token", "TKN", 18)
	data, err := EncodeERC20Approve(token, "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", big.NewInt(1000))
	if err != nil {
		t.Fatalf("EncodeERC20Approve() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeERC20Approve() returned empty data")
	}
}

func TestEncodeERC20TransferFrom(t *testing.T) {
	token := NewERC20Token("0xcontract", "Token", "TKN", 18)
	from := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	to := "0x1234567890123456789012345678901234567890"
	data, err := EncodeERC20TransferFrom(token, from, to, big.NewInt(500))
	if err != nil {
		t.Fatalf("EncodeERC20TransferFrom() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeERC20TransferFrom() returned empty data")
	}
}

func TestNewTokenTransfer(t *testing.T) {
	tp, err := NewTokenTransfer(
		"0xcontract",
		"0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		big.NewInt(1000),
		ChainMainnet,
	)
	if err != nil {
		t.Fatalf("NewTokenTransfer() error = %v", err)
	}
	if tp.Gas != GasLimitTokenTransfer.Uint64() {
		t.Errorf("Gas = %d, want %d", tp.Gas, GasLimitTokenTransfer.Uint64())
	}
}

func TestNewTokenApproval(t *testing.T) {
	tp, err := NewTokenApproval(
		"0xcontract",
		"0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		big.NewInt(1000),
		ChainMainnet,
	)
	if err != nil {
		t.Fatalf("NewTokenApproval() error = %v", err)
	}
	if tp.Gas != GasLimitTokenApproval.Uint64() {
		t.Errorf("Gas = %d, want %d", tp.Gas, GasLimitTokenApproval.Uint64())
	}
}

func TestNewERC721Token(t *testing.T) {
	token := NewERC721Token("0xcontract", "MyNFT", "NFT")
	if token == nil {
		t.Fatal("NewERC721Token() returned nil")
	}
}

func TestEncodeERC721Transfer(t *testing.T) {
	token := NewERC721Token("0xcontract", "MyNFT", "NFT")
	from := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	to := "0x1234567890123456789012345678901234567890"
	data, err := EncodeERC721Transfer(token, from, to, big.NewInt(1))
	if err != nil {
		t.Fatalf("EncodeERC721Transfer() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeERC721Transfer() returned empty data")
	}
}

func TestEncodeERC721Approve(t *testing.T) {
	token := NewERC721Token("0xcontract", "MyNFT", "NFT")
	approved := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	data, err := EncodeERC721Approve(token, approved, big.NewInt(1))
	if err != nil {
		t.Fatalf("EncodeERC721Approve() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeERC721Approve() returned empty data")
	}
}

func TestEncodeERC721SetApprovalForAll(t *testing.T) {
	token := NewERC721Token("0xcontract", "MyNFT", "NFT")
	operator := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	data, err := EncodeERC721SetApprovalForAll(token, operator, true)
	if err != nil {
		t.Fatalf("EncodeERC721SetApprovalForAll() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeERC721SetApprovalForAll() returned empty data")
	}
}

func TestCreateEventMonitor(t *testing.T) {
	monitor := CreateEventMonitor()
	if monitor == nil {
		t.Fatal("CreateEventMonitor() returned nil")
	}
}

func TestValidateAddress(t *testing.T) {
	valid := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	if !ValidateAddress(valid) {
		t.Errorf("ValidateAddress(%q) = false, want true", valid)
	}
}
