package edges

import (
	"encoding/json"
	"log"
	"math/big"
	"time"

	"dumb-api/config"
	"dumb-api/internal/graph"
	"dumb-api/internal/models"
	"dumb-api/internal/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EVMEdgeV2 struct {
	Token0     common.Address
	Token1     common.Address
	Reserve0   *big.Int
	Reserve1   *big.Int
	ZeroForOne bool
}

func (e *EVMEdgeV2) UpdateEdge(pendingLog types.Log, chainID string) {

	if !utils.HasTopics(pendingLog, config.SyncTopic) {
		return
	}

	data := utils.Chunks(common.Bytes2Hex(pendingLog.Data), 64)

	reserve0 := utils.StringToBigInt(data[0])
	reserve1 := utils.StringToBigInt(data[1])

	if e.ZeroForOne {
		e.Reserve0 = reserve0
		e.Reserve1 = reserve1
	} else {
		e.Reserve0 = reserve1
		e.Reserve1 = reserve0
	}

	amountOut := e.ComputeExactAmountOut(config.AmountIn)

	if amountOut == nil {
		log.Printf("[DELPHI] Nil AmountOut, Token0: %s, Token1: %s", e.Token0, e.Token1)
	}

	priceTick := &models.PriceTick{
		TokenIn:   e.Token0.String(),
		TokenOut:  e.Token1.String(),
		Chain:     chainID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := models.DB.Create(priceTick)
	if err != nil {
		log.Printf("[DELPHI] Err creating price tick V2: %v", err)
	}

}

func (e *EVMEdgeV2) ComputeExactAmountOut(amountIn *big.Int) *big.Int {

	amountInWithFee := new(big.Int).Mul(amountIn, big.NewInt(997))

	numerator := new(big.Int).Mul(amountInWithFee, e.Reserve1)
	denominator := new(big.Int).Add(new(big.Int).Mul(e.Reserve0, big.NewInt(1000)), amountInWithFee)

	// Compute amount out
	amountOut := new(big.Int).Quo(numerator, denominator)

	return amountOut
}
func (e *EVMEdgeV2) Export() []string {
	data := struct {
		Token0     string `json:"token0"`
		Token1     string `json:"token1"`
		Reserve0   string `json:"reserve0"`
		Reserve1   string `json:"reserve1"`
		ZeroForOne bool   `json:"zeroForOne"`
	}{
		Token0:     e.Token0.String(),
		Token1:     e.Token1.String(),
		Reserve0:   e.Reserve0.String(),
		Reserve1:   e.Reserve1.String(),
		ZeroForOne: e.ZeroForOne,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal edge data: %v", err)
		return []string{}
	}

	return []string{string(jsonData)}
}

func (e *EVMEdgeV2) ComputePriceImpact(amountIn *big.Int) *big.Float {

	return big.NewFloat(0)
}

func (e *EVMEdgeV2) Copy() graph.Edge {
	return &EVMEdgeV2{
		Token0:     e.Token0,
		Token1:     e.Token1,
		Reserve0:   e.Reserve0,
		Reserve1:   e.Reserve1,
		ZeroForOne: e.ZeroForOne,
	}
}

func (e *EVMEdgeV2) GetWeight() *big.Float {
	return big.NewFloat(0)
}
