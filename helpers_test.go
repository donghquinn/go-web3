package web3

import (
	"context"
	"math/big"
	"testing"

	blockchainhelper "github.com/donghquinn/go-blockchain-helper/pkg/web3"
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
	result, err := WeiToGwei(big.NewInt(1000000000))
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
	if result := FormatEther(wei, 18); result == "" {
		t.Error("FormatEther() returned empty string")
	}
}

func TestParseUnits(t *testing.T) {
	val, err := ParseUnits("1", 6)
	if err != nil {
		t.Fatalf("ParseUnits() error = %v", err)
	}
	if val.Cmp(big.NewInt(1000000)) != 0 {
		t.Errorf("ParseUnits(1, 6) = %v, want 1000000", val)
	}
}

func TestFormatUnits(t *testing.T) {
	if result := FormatUnits(big.NewInt(1000000), 6); result == "" {
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

func TestCreateTransactionWithEstimate(t *testing.T) {
	to := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	value, _ := EtherToWei("0.1")
	tp, err := CreateTransactionWithEstimate(to, value, nil, ChainMainnet)
	if err != nil {
		t.Fatalf("CreateTransactionWithEstimate() error = %v", err)
	}
	if tp.To != to {
		t.Errorf("To = %q, want %q", tp.To, to)
	}
	if tp.Gas == 0 {
		t.Error("Gas should be non-zero after estimation")
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
	gasPrice := big.NewInt(20000000000)
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

	price, err := GetOptimalGasPrice(context.Background(), NewClient(srv.URL), GasPriceFast)
	if err != nil {
		t.Fatalf("GetOptimalGasPrice() error = %v", err)
	}
	base, _ := ToWei("20", Gwei)
	if price.Cmp(base) <= 0 {
		t.Errorf("GetOptimalGasPrice(Fast) = %v should be greater than base %v", price, base)
	}
}

func TestGetOptimalGasPriceRPCError(t *testing.T) {
	srv := newMockRPCErrorServer(-32000, "internal error")
	defer srv.Close()

	_, err := GetOptimalGasPrice(context.Background(), NewClient(srv.URL), GasPriceFast)
	if err == nil {
		t.Fatal("expected RPC error")
	}
}

func TestEstimateGasWithBuffer(t *testing.T) {
	srv := newMockRPCServer("0x5208") // 21000
	defer srv.Close()

	gas, err := EstimateGasWithBuffer(context.Background(), NewClient(srv.URL), map[string]interface{}{
		"to": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
	}, 0.1)
	if err != nil {
		t.Fatalf("EstimateGasWithBuffer() error = %v", err)
	}
	if gas != 23100 { // 21000 + 21000*0.1
		t.Errorf("EstimateGasWithBuffer() = %d, want 23100", gas)
	}
}

func TestEstimateGasWithBufferRPCError(t *testing.T) {
	srv := newMockRPCErrorServer(-32000, "internal error")
	defer srv.Close()

	_, err := EstimateGasWithBuffer(context.Background(), NewClient(srv.URL), map[string]interface{}{}, 0.1)
	if err == nil {
		t.Fatal("expected RPC error")
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

func TestGetTokenBalance(t *testing.T) {
	// Server returns a 32-byte ABI-encoded uint256 (value = 42 = 0x2a)
	encoded := "0x000000000000000000000000000000000000000000000000000000000000002a"
	srv := newMockRPCServer(encoded)
	defer srv.Close()

	balance, err := GetTokenBalance(context.Background(), NewClient(srv.URL),
		"0xA0b86a33E6417c48cd7a94Ca95e70aD2c51e74f7",
		"0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
	)
	if err != nil {
		t.Fatalf("GetTokenBalance() error = %v", err)
	}
	if balance.Int64() != 42 {
		t.Errorf("GetTokenBalance() = %v, want 42", balance)
	}
}

func TestGetTokenBalanceRPCError(t *testing.T) {
	srv := newMockRPCErrorServer(-32000, "internal error")
	defer srv.Close()

	_, err := GetTokenBalance(context.Background(), NewClient(srv.URL), "0xcontract", "0xowner")
	if err == nil {
		t.Fatal("expected RPC error")
	}
}

func TestGetTokenAllowance(t *testing.T) {
	encoded := "0x0000000000000000000000000000000000000000000000000000000000000064" // 100
	srv := newMockRPCServer(encoded)
	defer srv.Close()

	allowance, err := GetTokenAllowance(context.Background(), NewClient(srv.URL),
		"0xA0b86a33E6417c48cd7a94Ca95e70aD2c51e74f7",
		"0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		"0x1234567890123456789012345678901234567890",
	)
	if err != nil {
		t.Fatalf("GetTokenAllowance() error = %v", err)
	}
	if allowance.Int64() != 100 {
		t.Errorf("GetTokenAllowance() = %v, want 100", allowance)
	}
}

func TestGetTokenBalanceFromHexError(t *testing.T) {
	// Server returns a result without 0x prefix → FromHex fails.
	srv := newMockRPCServer("deadbeef")
	defer srv.Close()

	_, err := GetTokenBalance(context.Background(), NewClient(srv.URL), "0xcontract", "0xowner")
	if err == nil {
		t.Fatal("expected FromHex error for result without 0x prefix")
	}
}

func TestGetTokenAllowanceRPCError(t *testing.T) {
	srv := newMockRPCErrorServer(-32000, "internal error")
	defer srv.Close()

	_, err := GetTokenAllowance(context.Background(), NewClient(srv.URL), "0xcontract", "0xowner", "0xspender")
	if err == nil {
		t.Fatal("expected RPC error")
	}
}

func TestGetTokenAllowanceFromHexError(t *testing.T) {
	srv := newMockRPCServer("deadbeef")
	defer srv.Close()

	_, err := GetTokenAllowance(context.Background(), NewClient(srv.URL), "0xcontract", "0xowner", "0xspender")
	if err == nil {
		t.Fatal("expected FromHex error for result without 0x prefix")
	}
}

func TestEncodeFunctionCallAdvanced(t *testing.T) {
	params := []blockchainhelper.ABIParam{
		{Type: "address"},
		{Type: "uint256"},
	}
	args := []interface{}{
		"0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		big.NewInt(1000),
	}
	data, err := EncodeFunctionCallAdvanced("transfer(address,uint256)", params, args)
	if err != nil {
		t.Fatalf("EncodeFunctionCallAdvanced() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeFunctionCallAdvanced() returned empty data")
	}
}

func TestDecodeFunctionResult(t *testing.T) {
	// ABI-encode a uint256 = 42, then decode it back
	encoded, err := EncodeFunctionCallAdvanced("getValue()", nil, nil)
	_ = encoded
	// Decode a raw ABI uint256 (32 bytes) representing 42
	raw := make([]byte, 32)
	raw[31] = 42
	result, err := DecodeFunctionResult([]string{"uint256"}, raw)
	if err != nil {
		t.Fatalf("DecodeFunctionResult() error = %v", err)
	}
	if len(result) == 0 {
		t.Error("DecodeFunctionResult() returned empty result")
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
	data, err := EncodeERC721Approve(token, "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", big.NewInt(1))
	if err != nil {
		t.Fatalf("EncodeERC721Approve() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeERC721Approve() returned empty data")
	}
}

func TestEncodeERC721SetApprovalForAll(t *testing.T) {
	token := NewERC721Token("0xcontract", "MyNFT", "NFT")
	data, err := EncodeERC721SetApprovalForAll(token, "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", true)
	if err != nil {
		t.Fatalf("EncodeERC721SetApprovalForAll() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("EncodeERC721SetApprovalForAll() returned empty data")
	}
}

func TestCreateEventMonitor(t *testing.T) {
	if monitor := CreateEventMonitor(); monitor == nil {
		t.Fatal("CreateEventMonitor() returned nil")
	}
}

func TestParseTransferEvent(t *testing.T) {
	// Construct a minimal Transfer event with the standard 3 topics and 32-byte data.
	from := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	to := "0x1234567890123456789012345678901234567890"
	// Transfer(address indexed from, address indexed to, uint256 value)
	transferSig := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	padAddr := func(addr string) string {
		// pad address to 32 bytes (64 hex chars) with leading zeros
		return "0x000000000000000000000000" + addr[2:]
	}
	// value = 100 in 32-byte hex
	value := "0x0000000000000000000000000000000000000000000000000000000000000064"

	event := blockchainhelper.Event{
		Address: "0xA0b86a33E6417c48cd7a94Ca95e70aD2c51e74f7",
		Topics:  []string{transferSig, padAddr(from), padAddr(to)},
		Data:    value,
	}

	result, err := ParseTransferEvent(event)
	if err != nil {
		t.Fatalf("ParseTransferEvent() error = %v", err)
	}
	if result == nil {
		t.Fatal("ParseTransferEvent() returned nil")
	}
}

func TestValidateAddress(t *testing.T) {
	valid := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	if !ValidateAddress(valid) {
		t.Errorf("ValidateAddress(%q) = false, want true", valid)
	}
}

func TestPrivateKeyToAddressHelper(t *testing.T) {
	const privHex = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	addr, err := PrivateKeyToAddressHelper(privHex)
	if err != nil {
		t.Fatalf("PrivateKeyToAddressHelper() error = %v", err)
	}
	if !IsAddress(addr) {
		t.Errorf("PrivateKeyToAddressHelper() = %q is not a valid address", addr)
	}
}
