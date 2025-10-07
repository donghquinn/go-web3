package web3

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func ToWei(value string, unit string) (*big.Int, error) {
	val, ok := new(big.Float).SetString(value)
	if !ok {
		return nil, fmt.Errorf("invalid value: %s", value)
	}

	var multiplier *big.Int
	switch strings.ToLower(unit) {
	case "wei":
		multiplier = big.NewInt(1)
	case "kwei", "babbage", "femtoether":
		multiplier = big.NewInt(1e3)
	case "mwei", "lovelace", "picoether":
		multiplier = big.NewInt(1e6)
	case "gwei", "shannon", "nanoether", "nano":
		multiplier = big.NewInt(1e9)
	case "szabo", "microether", "micro":
		multiplier = big.NewInt(1e12)
	case "finney", "milliether", "milli":
		multiplier = big.NewInt(1e15)
	case "ether", "eth":
		multiplier = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	case "kether", "grand":
		multiplier = new(big.Int).Exp(big.NewInt(10), big.NewInt(21), nil)
	case "mether":
		multiplier = new(big.Int).Exp(big.NewInt(10), big.NewInt(24), nil)
	case "gether":
		multiplier = new(big.Int).Exp(big.NewInt(10), big.NewInt(27), nil)
	case "tether":
		multiplier = new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil)
	default:
		return nil, fmt.Errorf("unknown unit: %s", unit)
	}

	result := new(big.Float).Mul(val, new(big.Float).SetInt(multiplier))
	wei, _ := result.Int(nil)
	return wei, nil
}

func FromWei(wei *big.Int, unit string) (string, error) {
	if wei == nil {
		return "0", nil
	}

	var divisor *big.Int
	switch strings.ToLower(unit) {
	case "wei":
		divisor = big.NewInt(1)
	case "kwei", "babbage", "femtoether":
		divisor = big.NewInt(1e3)
	case "mwei", "lovelace", "picoether":
		divisor = big.NewInt(1e6)
	case "gwei", "shannon", "nanoether", "nano":
		divisor = big.NewInt(1e9)
	case "szabo", "microether", "micro":
		divisor = big.NewInt(1e12)
	case "finney", "milliether", "milli":
		divisor = big.NewInt(1e15)
	case "ether", "eth":
		divisor = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	case "kether", "grand":
		divisor = new(big.Int).Exp(big.NewInt(10), big.NewInt(21), nil)
	case "mether":
		divisor = new(big.Int).Exp(big.NewInt(10), big.NewInt(24), nil)
	case "gether":
		divisor = new(big.Int).Exp(big.NewInt(10), big.NewInt(27), nil)
	case "tether":
		divisor = new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil)
	default:
		return "", fmt.Errorf("unknown unit: %s", unit)
	}

	result := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(divisor))
	return result.String(), nil
}

func IsAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	if len(address) != 42 {
		return false
	}
	
	for _, r := range address[2:] {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return true
}

func ToHex(value interface{}) string {
	switch v := value.(type) {
	case int:
		return fmt.Sprintf("0x%x", v)
	case int64:
		return fmt.Sprintf("0x%x", v)
	case uint64:
		return fmt.Sprintf("0x%x", v)
	case *big.Int:
		return fmt.Sprintf("0x%x", v)
	case []byte:
		return fmt.Sprintf("0x%x", v)
	case string:
		if strings.HasPrefix(v, "0x") {
			return v
		}
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return fmt.Sprintf("0x%x", val)
		}
		return fmt.Sprintf("0x%x", []byte(v))
	default:
		return fmt.Sprintf("0x%x", fmt.Sprintf("%v", v))
	}
}

func FromHex(hex string) (*big.Int, error) {
	if !strings.HasPrefix(hex, "0x") {
		return nil, fmt.Errorf("hex string must start with 0x")
	}
	
	value := new(big.Int)
	value.SetString(hex[2:], 16)
	return value, nil
}

func PadLeft(str string, length int, padChar string) string {
	for len(str) < length {
		str = padChar + str
	}
	return str
}

func PadRight(str string, length int, padChar string) string {
	for len(str) < length {
		str = str + padChar
	}
	return str
}