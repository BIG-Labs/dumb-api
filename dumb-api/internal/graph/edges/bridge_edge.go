package edges

import (
	"encoding/json"
	"log"
	"math/big"

	"dumb-api/internal/graph"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BridgeEdge struct {
	Token0 common.Address
	Token1 common.Address
}

func (e *BridgeEdge) UpdateEdge(pendingLog types.Log, chainID string) {

}

func (e *BridgeEdge) ComputeExactAmountOut(amountIn *big.Int) *big.Int {
	return amountIn
}

func (e *BridgeEdge) ComputePriceImpact(amountIn *big.Int) *big.Float {

	return big.NewFloat(0)
}

func (e *BridgeEdge) Export() []string {
	data := struct {
		Token0 string `json:"token0"`
		Token1 string `json:"token1"`
	}{
		Token0: e.Token0.String(),
		Token1: e.Token1.String(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal edge data: %v", err)
		return []string{}
	}

	return []string{string(jsonData)}
}

func (e *BridgeEdge) Copy() graph.Edge {
	return &BridgeEdge{
		Token0: e.Token0,
		Token1: e.Token1,
	}
}

func (e *BridgeEdge) GetWeight() *big.Float {
	return big.NewFloat(0)
}
