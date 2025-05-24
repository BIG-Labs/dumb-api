package graph

import (
	"math/big"
	"github.com/ethereum/go-ethereum/core/types"
)

type Edge interface {
	UpdateEdge(vLog types.Log, chainID string)
	ComputeExactAmountOut(amountIn *big.Int) *big.Int
	ComputePriceImpact(amountIn *big.Int) *big.Float
	Export() []string
	Copy() Edge
	GetWeight() *big.Float
}
