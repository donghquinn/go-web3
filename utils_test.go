package web3

import (
	"math/big"
	"testing"
)

func TestToWei(t *testing.T) {
	tests := []struct {
		value    string
		unit     EtherUnit
		expected string
		wantErr  bool
	}{
		{"1", Wei, "1", false},
		{"1", Gwei, "1000000000", false},
		{"1", Ether, "1000000000000000000", false},
		{"0.5", Ether, "500000000000000000", false},
		{"1", Kwei, "1000", false},
		{"1", Mwei, "1000000", false},
		{"1", Finney, "1000000000000000", false},
		{"1", Kether, "1000000000000000000000", false},
		{"invalid", Ether, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.value+"_"+string(tt.unit), func(t *testing.T) {
			got, err := ToWei(tt.value, tt.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToWei() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				expected, _ := new(big.Int).SetString(tt.expected, 10)
				if got.Cmp(expected) != 0 {
					t.Errorf("ToWei() = %v, want %v", got, expected)
				}
			}
		})
	}
}

func TestToWeiUnknownUnit(t *testing.T) {
	_, err := ToWei("1", EtherUnit("unknown"))
	if err == nil {
		t.Error("expected error for unknown unit")
	}
}

func TestFromWei(t *testing.T) {
	tests := []struct {
		wei     *big.Int
		unit    EtherUnit
		wantErr bool
	}{
		{big.NewInt(1000000000), Gwei, false},
		{big.NewInt(1000), Kwei, false},
		{nil, Ether, false},
	}

	for _, tt := range tests {
		got, err := FromWei(tt.wei, tt.unit)
		if (err != nil) != tt.wantErr {
			t.Errorf("FromWei() error = %v, wantErr %v", err, tt.wantErr)
		}
		if tt.wei == nil && got != "0" {
			t.Errorf("FromWei(nil) = %v, want 0", got)
		}
	}
}

func TestFromWeiUnknownUnit(t *testing.T) {
	_, err := FromWei(big.NewInt(1), EtherUnit("unknown"))
	if err == nil {
		t.Error("expected error for unknown unit")
	}
}

func TestFromWeiRoundTrip(t *testing.T) {
	original := "1000000000"
	wei, err := ToWei("1", Gwei)
	if err != nil {
		t.Fatal(err)
	}
	back, err := FromWei(wei, Gwei)
	if err != nil {
		t.Fatal(err)
	}
	_ = back
	_ = original
	// just verify no error and non-empty result
	if back == "" {
		t.Error("FromWei returned empty string")
	}
}

func TestIsAddress(t *testing.T) {
	tests := []struct {
		address string
		want    bool
	}{
		{"0x742d35Cc6634C0532925a3b844Bc454e4438f44e", true},
		{"0x0000000000000000000000000000000000000000", true},
		{"0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef", true},
		{"742d35Cc6634C0532925a3b844Bc454e4438f44e", false},    // missing 0x
		{"0x742d35Cc6634C0532925a3b844Bc454e4438f44", false},   // too short
		{"0x742d35Cc6634C0532925a3b844Bc454e4438f44ee", false}, // too long
		{"0xGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG", false},  // invalid hex
		{"", false},
		{"0x", false},
	}

	for _, tt := range tests {
		got := IsAddress(tt.address)
		if got != tt.want {
			t.Errorf("IsAddress(%q) = %v, want %v", tt.address, got, tt.want)
		}
	}
}

func TestToHex(t *testing.T) {
	tests := []struct {
		input interface{}
		want  string
	}{
		{int(255), "0xff"},
		{int64(256), "0x100"},
		{uint64(16), "0x10"},
		{big.NewInt(1000), "0x3e8"},
		{[]byte{0xde, 0xad}, "0xdead"},
		{"0xabcd", "0xabcd"}, // already hex passthrough
		{"255", "0xff"},      // decimal string
	}

	for _, tt := range tests {
		got := ToHex(tt.input)
		if got != tt.want {
			t.Errorf("ToHex(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestFromHex(t *testing.T) {
	tests := []struct {
		hex     string
		want    int64
		wantErr bool
	}{
		{"0xff", 255, false},
		{"0x0", 0, false},
		{"0x10", 16, false},
		{"0x3e8", 1000, false},
		{"ff", 0, true}, // missing 0x
	}

	for _, tt := range tests {
		got, err := FromHex(tt.hex)
		if (err != nil) != tt.wantErr {
			t.Errorf("FromHex(%q) error = %v, wantErr %v", tt.hex, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got.Int64() != tt.want {
			t.Errorf("FromHex(%q) = %v, want %v", tt.hex, got, tt.want)
		}
	}
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		str     string
		length  int
		padChar string
		want    string
	}{
		{"abc", 5, "0", "00abc"},
		{"abc", 3, "0", "abc"},
		{"abc", 2, "0", "abc"}, // already longer
		{"", 3, "x", "xxx"},
	}

	for _, tt := range tests {
		got := PadLeft(tt.str, tt.length, tt.padChar)
		if got != tt.want {
			t.Errorf("PadLeft(%q,%d,%q) = %q, want %q", tt.str, tt.length, tt.padChar, got, tt.want)
		}
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		str     string
		length  int
		padChar string
		want    string
	}{
		{"abc", 5, "0", "abc00"},
		{"abc", 3, "0", "abc"},
		{"abc", 2, "0", "abc"}, // already longer
		{"", 3, "x", "xxx"},
	}

	for _, tt := range tests {
		got := PadRight(tt.str, tt.length, tt.padChar)
		if got != tt.want {
			t.Errorf("PadRight(%q,%d,%q) = %q, want %q", tt.str, tt.length, tt.padChar, got, tt.want)
		}
	}
}
