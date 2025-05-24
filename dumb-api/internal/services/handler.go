package services

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gobuffalo/pop/v6"
)

type EVMHandler interface {
	UpdateLastBlock(db *pop.Connection, block int64) error
	LastBlock(ctx context.Context, db *pop.Connection) (int64, error)
	Block(ctx context.Context, height *int64) (*types.Block, error)
	BlockResult(ctx context.Context, height *int64) (*types.Block, error)
	Listen(db *pop.Connection, block *types.Block) error
	GetChainName() string
	GetChainID() string
}
