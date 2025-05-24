package config

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
)

var (
	SyncTopic         string
	SwapV3Topic       string
	AmountIn          *big.Int
	DefaultBuilderFee *big.Int
	FeeTiers          []string
	TokensByChain     map[string][]string
	EVMConfig         map[string]ChainConfig
	SavePath          string
	AVALANCHE_RPC_URL string
	ENV               string
)

type TokenConfig struct {
	Address     string `json:"address"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	TokenHome   string `json:"token_home"`
	TokenRemote string `json:"token_remote"`
}

type DexConfig struct {
	Factories []string `json:"factories"`
}

type ChainConfig struct {
	UniswapV2 DexConfig     `json:"UniswapV2"`
	UniswapV3 DexConfig     `json:"UniswapV3"`
	Tokens    []TokenConfig `json:"tokens"`
	ChainId   int           `json:"chainId"`
}

func init() {
	var err error

	TokensByChain = make(map[string][]string)
	EVMConfig = make(map[string]ChainConfig)

	SyncTopic = os.Getenv("SYNC_TOPIC")
	SwapV3Topic = os.Getenv("SWAP_V3_TOPIC")
	AVALANCHE_RPC_URL = os.Getenv("AVALANCHE_RPC_URL")
	ENV = os.Getenv("ENV")

	AmountIn, err = loadBigInt("AMOUNT_IN")
	if err != nil {
		log.Fatalf("failed to load AMOUNT_IN: %v", err)
	}

	DefaultBuilderFee, err = loadBigInt("DEFAULT_BUILDER_FEE")
	if err != nil {
		log.Fatalf("failed to load DEFAULT_BUILDER_FEE: %v", err)
	}

	FeeTiers, err = loadFeeTiers()
	if err != nil {
		log.Fatalf("failed to load FEE_TIERS: %v", err)
	}

	if err := loadEVMConfig(); err != nil {
		log.Fatalf("failed to load EVM config: %v", err)
	}

	SavePath = "data/backup"
}

func loadBigInt(envVar string) (*big.Int, error) {
	value := os.Getenv(envVar)
	if value == "" {
		return nil, fmt.Errorf("%s environment variable not set", envVar)
	}
	result := new(big.Int)
	if _, ok := result.SetString(value, 10); !ok {
		return nil, fmt.Errorf("invalid %s value: %s", envVar, value)
	}
	return result, nil
}

func loadFeeTiers() ([]string, error) {
	tiers := os.Getenv("FEE_TIERS")
	if tiers == "" {
		return nil, fmt.Errorf("FEE_TIERS environment variable not set")
	}
	return strings.Split(tiers, ","), nil
}

func loadEVMConfig() error {
	jsonData, err := os.ReadFile("evm_config.json")
	if err != nil {
		return fmt.Errorf("error reading evm_config.json: %w", err)
	}

	if err := json.Unmarshal(jsonData, &EVMConfig); err != nil {
		return fmt.Errorf("error parsing evm_config.json: %w", err)
	}

	for chain, chainConfig := range EVMConfig {
		addresses := make([]string, len(chainConfig.Tokens))
		for i, token := range chainConfig.Tokens {
			addresses[i] = token.Address
		}
		TokensByChain[chain] = addresses
	}

	return nil
}
