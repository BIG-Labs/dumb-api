package utils

import (
	"log"
	"strings"

	"dumb-api/config"
)

// GetTokenMappings returns the home and remote token addresses for a given token on a specific chain
func GetTokenMappings(token, chain string) (string, string) {
	log.Printf("Getting token mappings for token %s on chain %s", token, chain)

	// Convert chain name to uppercase to match config
	chainUpper := strings.ToUpper(chain)
	chainConfig, exists := config.EVMConfig[chainUpper]
	if !exists {
		log.Printf("Chain %s (as %s) not found in EVMConfig", chain, chainUpper)
		return "", ""
	}

	tokenLower := strings.ToLower(token)
	log.Printf("Token lower: %s", tokenLower)
	for _, tokenConfig := range chainConfig.Tokens {
		configAddrLower := strings.ToLower(tokenConfig.Address)
		log.Printf("Comparing token %s with config token %s", tokenLower, configAddrLower)
		if configAddrLower == tokenLower {
			log.Printf("Found matching token. Home: %s, Remote: %s", tokenConfig.TokenHome, tokenConfig.TokenRemote)
			return tokenConfig.TokenHome, tokenConfig.TokenRemote
		}
	}

	log.Printf("No matching token found for %s in chain %s", tokenLower, chainUpper)
	return "", ""
}
