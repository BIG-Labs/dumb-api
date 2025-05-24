package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/0x7183/unifi-backend/internal/utils"
)

type GlacierTokenResponse struct {
	Address           string                 `json:"address"`
	DeploymentDetails *DeploymentDetails     `json:"deploymentDetails"`
	Name              string                 `json:"name"`
	LogoAsset         *LogoAsset             `json:"logoAsset"`
	Color             string                 `json:"color"`
	ResourceLinks     []any        `json:"resourceLinks"`
	ErcType           string                 `json:"ercType"`
	Symbol            string                 `json:"symbol"`
	Decimals          int                    `json:"decimals"`
	PricingProviders  map[string]any `json:"pricingProviders"`
}

type DeploymentDetails struct {
	TxHash          string `json:"txHash"`
	DeployerAddress string `json:"deployerAddress"`
}

type LogoAsset struct {
	ImageUri string `json:"imageUri"`
}

type GlacierAPIService struct {
	client *http.Client
}

func NewGlacierAPIService() *GlacierAPIService {
	return &GlacierAPIService{
		client: utils.NewHTTPClient(),
	}
}

func (g *GlacierAPIService) GetTokenInfo(chainId int, tokenAddress string) (*GlacierTokenResponse, error) {
	url := fmt.Sprintf("https://glacier-api.avax.network/v1/chains/%d/addresses/%s", chainId, tokenAddress)

	resp, err := g.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Glacier API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get token info from Glacier API: %d", resp.StatusCode)
	}

	var tokenInfo GlacierTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenInfo, nil
}
