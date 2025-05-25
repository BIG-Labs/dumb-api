package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"dumb-api/config"
	"dumb-api/internal/graph"
	"dumb-api/internal/models"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
)

type AvalancheHandler struct {
	Client  *ethclient.Client
	ChainID string
}

func NewAvalancheHandler(rpcURL string) *AvalancheHandler {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to create EVM client: %v", err)
	}
	return &AvalancheHandler{
		Client:  client,
		ChainID: "43114",
	}
}

func (h *AvalancheHandler) UpdateLastBlock(db *pop.Connection, block int64) error {
	var chainState models.ChainState

	err := db.Where("chain_id = ?", h.ChainID).First(&chainState)
	if err != nil {
		chainState = models.ChainState{
			ID:        uuid.Must(uuid.NewV4()),
			ChainID:   h.ChainID,
			LastBlock: block,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		verrs, err := db.ValidateAndCreate(&chainState)
		if err != nil {
			return fmt.Errorf("failed to create chain state: %v", err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("validation errors creating chain state: %v", verrs.Errors)
		}
		return nil
	}

	chainState.LastBlock = block
	chainState.UpdatedAt = time.Now()

	verrs, err := db.ValidateAndUpdate(&chainState)
	if err != nil {
		return fmt.Errorf("failed to update chain state: %v", err)
	}
	if verrs.HasAny() {
		return fmt.Errorf("validation errors updating chain state: %v", verrs.Errors)
	}

	return nil
}

func (h *AvalancheHandler) LastBlock(ctx context.Context, db *pop.Connection) (int64, error) {
	var chainState models.ChainState

	err := db.Where("chain_id = ?", h.ChainID).First(&chainState)
	if err != nil {
		return 0, nil
	}

	return chainState.LastBlock, nil
}

func (h *AvalancheHandler) Block(ctx context.Context, height *int64) (*types.Block, error) {
	var blockNumber *big.Int
	if height != nil {
		blockNumber = big.NewInt(*height)
	}

	block, err := h.Client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %v", err)
	}

	return block, nil
}

func (h *AvalancheHandler) BlockResult(ctx context.Context, height *int64) (*types.Block, error) {
	return h.Block(ctx, height)
}

func (h *AvalancheHandler) Listen(db *pop.Connection, block *types.Block) error {
	for _, tx := range block.Transactions() {
		receipt, err := h.Client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Printf("Failed to get transaction receipt: %v", err)
			continue
		}

		for _, vLog := range receipt.Logs {
			if len(vLog.Topics) == 0 {
				continue
			}

			topic := vLog.Topics[0].Hex()
			if topic != config.SyncTopic || topic != config.SwapV3Topic {
				continue
			}

			poolAddr := vLog.Address.Hex()
			pool := graph.GetGlobalGraph().GetPool(poolAddr)
			if pool == nil {
				log.Printf("Pool not found for address: %s", poolAddr)
				continue
			}

			edge01 := graph.GetGlobalGraph().GetEdge(pool.Token0, pool.Token1, pool.Pair, h.ChainID)
			edge10 := graph.GetGlobalGraph().GetEdge(pool.Token1, pool.Token0, pool.Pair, h.ChainID)

			if edge01 != nil {
				edge01.UpdateEdge(*vLog, h.ChainID)
			}

			if edge10 != nil {
				edge10.UpdateEdge(*vLog, h.ChainID)
			}
		}
	}
	return nil
}

func (h *AvalancheHandler) GetChainName() string {
	return "Avalanche"
}

func (h *AvalancheHandler) GetChainID() string {
	return h.ChainID
}
