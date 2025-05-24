package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/0x7183/unifi-backend/internal/utils"
)

type CoinGeckoService struct {
	client *http.Client
	apiKey string
}

type TokenPriceResponse map[string]struct {
	USD float64 `json:"usd"`
}

func NewCoinGeckoService() *CoinGeckoService {
	return &CoinGeckoService{
		client: &http.Client{},
		apiKey: "CG-ZcF22TjagYiGGXmtsA1LnmbX",
	}
}

func (s *CoinGeckoService) GetTokenPrice(chainName, contractAddress string) (float64, error) {
	chainName = strings.ToLower(chainName)

	// Get the mapped address for CoinGecko
	mappedAddress := utils.MapCoinGeckoAddress(contractAddress)
	mappedAddress = strings.ToLower(mappedAddress)

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/token_price/%s?contract_addresses=%s&vs_currencies=usd&x_cg_demo_api_key=%s",
		chainName, mappedAddress, s.apiKey)

	fmt.Printf("CoinGecko URL: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var priceResponse TokenPriceResponse
	if err := json.Unmarshal(body, &priceResponse); err != nil {
		return 0, fmt.Errorf("failed to parse response: %v", err)
	}

	// Look for the price using the mapped address
	price, ok := priceResponse[mappedAddress]
	if !ok {
		return 0, fmt.Errorf("price not found for token %s (mapped from %s)", mappedAddress, contractAddress)
	}

	return price.USD, nil
}
