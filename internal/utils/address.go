package utils

import "strings"

var addressMapping = map[string]string{
	"0x408d4cd0adb7cebd1f1a1c33a0ba2098e1295bab": "0x152b9d0fdc40c096757f570a51e494bd4b943e50", // WBTC
}

func MapCoinGeckoAddress(address string) string {
	lowerAddress := strings.ToLower(address)

	if mapped, exists := addressMapping[lowerAddress]; exists {
		return mapped
	}

	return address
}


