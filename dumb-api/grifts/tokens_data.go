package grifts

import (
	"fmt"
	"log"
	"time"

	"dumb-api/config"
	"dumb-api/internal/services"
	"dumb-api/models"

	"github.com/gobuffalo/grift/grift"
	"github.com/gofrs/uuid"
)

var _ = grift.Namespace("tokens", func() {
	grift.Desc("seed", "Seeds or updates the tokens data")
	grift.Add("seed", func(c *grift.Context) error {
		log.Println("Starting tokens data update...")

		glacierService := services.NewGlacierAPIService()
		coingeckoService := services.NewCoinGeckoService()
		now := time.Now()

		for chainName, chainConfig := range config.EVMConfig {
			log.Printf("Processing chain: %s", chainName)

			for _, tokenConfig := range chainConfig.Tokens {
				log.Printf("Processing token: %s (%s)", tokenConfig.Symbol, tokenConfig.Address)

				if chainName == "COQNET" && tokenConfig.Symbol == "USDC" {
					token := models.Token{
						ID:        uuid.Must(uuid.NewV4()),
						Address:   tokenConfig.Address,
						ChainID:   fmt.Sprintf("%d", chainConfig.ChainId),
						Icon:      "",
						Name:      tokenConfig.Name,
						Symbol:    tokenConfig.Symbol,
						Price:     1.0,
						Decimals:  6,
						UpdatedAt: now,
					}

					err := models.DB.Create(&token)
					if err != nil {
						log.Printf("Failed to create token record %s: %v", tokenConfig.Address, err)
						continue
					}
					log.Printf("Successfully created new token price record: %s (%s) with price %v at %v", token.Symbol, token.Name, token.Price, now)
					continue
				}

				tokenInfo, err := glacierService.GetTokenInfo(chainConfig.ChainId, tokenConfig.Address)
				if err != nil {
					log.Printf("Failed to fetch token info for %s: %v", tokenConfig.Address, err)
					continue
				}

				var iconURL string
				if tokenInfo.LogoAsset != nil {
					iconURL = tokenInfo.LogoAsset.ImageUri
				}

				price, err := coingeckoService.GetTokenPrice(chainName, tokenConfig.Address)
				if err != nil {
					log.Printf("Failed to fetch token price for %s: %v", tokenConfig.Address, err)
					continue
				}

				// Create a new token record for each price update
				token := models.Token{
					ID:        uuid.Must(uuid.NewV4()),
					Address:   tokenConfig.Address,
					ChainID:   fmt.Sprintf("%d", chainConfig.ChainId),
					Icon:      iconURL,
					Name:      tokenConfig.Name,
					Symbol:    tokenConfig.Symbol,
					Price:     price,
					Decimals:  tokenInfo.Decimals,
					UpdatedAt: now,
				}

				// Create a new record for each price update
				err = models.DB.Create(&token)
				if err != nil {
					log.Printf("Failed to create token record %s: %v", tokenConfig.Address, err)
					continue
				}
				log.Printf("Successfully created new token price record: %s (%s) with price %v at %v", token.Symbol, token.Name, token.Price, now)
			}
		}

		log.Println("Tokens data update completed!")
		return nil
	})
})
