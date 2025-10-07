package web3

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func ToWei(value string, unit EtherUnit) (*big.Int, error) {
	val, ok := new(big.Float).SetString(value)
	if !ok {
		return nil, fmt.Errorf("invalid value: %s", value)
	}

	var multiplier *big.Int
	switch unit {
	case Wei:
		multiplier = big.NewInt(1)
	case Kwei, Babbage, Femtoether:
		multiplier = big.NewInt(1e3)
	case Mwei, Lovelace, Picoether:
		multiplier = big.NewInt(1e6)
	case Gwei, Shannon, Nanoether, Nano:
		multiplier = big.NewInt(1e9)
	case Szabo, Microether, Micro:
		multiplier = big.NewInt(1e12)
	case Finney, Milliether, Milli:
		multiplier = big.NewInt(1e15)
	case Ether, EthUnit:
		multiplier = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	case Kether, Grand:
		multiplier = new(big.Int).Exp(big.NewInt(10), big.NewInt(21), nil)
	case Mether:
		multiplier = new(big.Int).Exp(big.NewInt(10), big.NewInt(24), nil)
	case Gether:
		multiplier = new(big.Int).Exp(big.NewInt(10), big.NewInt(27), nil)
	case Tether:
		multiplier = new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil)
	default:
		return nil, fmt.Errorf("unknown unit: %s", unit)
	}

	result := new(big.Float).Mul(val, new(big.Float).SetInt(multiplier))
	wei, _ := result.Int(nil)
	return wei, nil
}

func FromWei(wei *big.Int, unit EtherUnit) (string, error) {
	if wei == nil {
		return "0", nil
	}

	var divisor *big.Int
	switch unit {
	case Wei:
		divisor = big.NewInt(1)
	case Kwei, Babbage, Femtoether:
		divisor = big.NewInt(1e3)
	case Mwei, Lovelace, Picoether:
		divisor = big.NewInt(1e6)
	case Gwei, Shannon, Nanoether, Nano:
		divisor = big.NewInt(1e9)
	case Szabo, Microether, Micro:
		divisor = big.NewInt(1e12)
	case Finney, Milliether, Milli:
		divisor = big.NewInt(1e15)
	case Ether, EthUnit:
		divisor = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	case Kether, Grand:
		divisor = new(big.Int).Exp(big.NewInt(10), big.NewInt(21), nil)
	case Mether:
		divisor = new(big.Int).Exp(big.NewInt(10), big.NewInt(24), nil)
	case Gether:
		divisor = new(big.Int).Exp(big.NewInt(10), big.NewInt(27), nil)
	case Tether:
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