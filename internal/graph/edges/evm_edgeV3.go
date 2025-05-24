package edges

import (
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"time"

	"github.com/0x7183/unifi-backend/config"
	"github.com/0x7183/unifi-backend/internal/graph"
	"github.com/0x7183/unifi-backend/internal/models"
	"github.com/0x7183/unifi-backend/internal/utils"
	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/daoleno/uniswapv3-sdk/entities"
	uniswapv3utils "github.com/daoleno/uniswapv3-sdk/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrSqrtPriceLimitX96TooLow  = errors.New("SqrtPriceLimitX96 too low")
	ErrSqrtPriceLimitX96TooHigh = errors.New("SqrtPriceLimitX96 too high")
)

type StepComputations struct {
	sqrtPriceStartX96 *big.Int
	tickNext          int
	initialized       bool
	sqrtPriceNextX96  *big.Int
	amountIn          *big.Int
	amountOut         *big.Int
	feeAmount         *big.Int
}

type EVMEdgeV3 struct {
	Token0           common.Address
	Token1           common.Address
	Fee              constants.FeeAmount
	SqrtRatioX96     *big.Int
	Liquidity        *big.Int
	TickCurrent      int
	TickDataProvider entities.TickDataProvider
	ZeroForOne       bool
	exchangeRate     *big.Float
}

func (e *EVMEdgeV3) Copy() graph.Edge {
	return &EVMEdgeV3{
		Token0:           e.Token0,
		Token1:           e.Token1,
		Fee:              e.Fee,
		SqrtRatioX96:     e.SqrtRatioX96,
		Liquidity:        e.Liquidity,
		TickCurrent:      e.TickCurrent,
		TickDataProvider: e.TickDataProvider,
		ZeroForOne:       e.ZeroForOne,
		exchangeRate:     e.exchangeRate,
	}
}

func (e *EVMEdgeV3) UpdateEdge(pendingLog types.Log, chainID string) {
	if !utils.HasTopics(pendingLog, config.SwapV3Topic) {
		return
	}

	data := utils.Chunks(common.Bytes2Hex(pendingLog.Data), 64)
	sqrtRatioX96 := utils.StringToBigInt(data[2])
	liquidity := utils.StringToBigInt(data[3])
	tickCurrent := int(utils.StringToBigInt(data[4]).Int64())

	e.TickCurrent = tickCurrent
	e.SqrtRatioX96 = sqrtRatioX96
	e.Liquidity = liquidity

	amountOut := e.ComputeExactAmountOut(config.AmountIn)

	if amountOut == nil {
		amountOut = constants.Zero
	}

	e.exchangeRate = new(big.Float).Quo(new(big.Float).SetInt(amountOut), new(big.Float).SetInt(config.AmountIn))

	priceTick := &models.PriceTick{
		TokenIn:   e.Token0.String(),
		TokenOut:  e.Token1.String(),
		Chain:     chainID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := models.DB.Create(priceTick)
	if err != nil {
		log.Printf("[DELPHI] Err creating price tick V3: %v", err)
	}

}

func (e *EVMEdgeV3) ComputeExactAmountOut(inputAmount *big.Int) *big.Int {
	zeroForOne := e.ZeroForOne
	outputAmount, _, _, _, err := e.swap(zeroForOne, inputAmount, nil)
	if err != nil {
		log.Printf("[DELPHI] Err calculating amount out V3: %v", err)
		return new(big.Int).Set(constants.Zero)
	}

	return new(big.Int).Mul(outputAmount, constants.NegativeOne)
}

func (e *EVMEdgeV3) swap(zeroForOne bool, amountSpecified, sqrtPriceLimitX96 *big.Int) (amountCalculated *big.Int, sqrtRatioX96 *big.Int, liquidity *big.Int, tickCurrent int, err error) {
	if sqrtPriceLimitX96 == nil {
		if zeroForOne {
			sqrtPriceLimitX96 = new(big.Int).Add(uniswapv3utils.MinSqrtRatio, constants.One)
		} else {
			sqrtPriceLimitX96 = new(big.Int).Sub(uniswapv3utils.MaxSqrtRatio, constants.One)
		}
	}

	if zeroForOne {
		if sqrtPriceLimitX96.Cmp(uniswapv3utils.MinSqrtRatio) <= 0 {
			return nil, nil, nil, 0, ErrSqrtPriceLimitX96TooLow
		}
		if sqrtPriceLimitX96.Cmp(e.SqrtRatioX96) >= 0 {
			return nil, nil, nil, 0, ErrSqrtPriceLimitX96TooHigh
		}
	} else {
		if sqrtPriceLimitX96.Cmp(uniswapv3utils.MaxSqrtRatio) >= 0 {
			return nil, nil, nil, 0, ErrSqrtPriceLimitX96TooHigh
		}
		if sqrtPriceLimitX96.Cmp(e.SqrtRatioX96) <= 0 {
			return nil, nil, nil, 0, ErrSqrtPriceLimitX96TooLow
		}
	}

	exactInput := amountSpecified.Cmp(constants.Zero) >= 0

	state := struct {
		amountSpecifiedRemaining *big.Int
		amountCalculated         *big.Int
		sqrtPriceX96             *big.Int
		tick                     int
		liquidity                *big.Int
	}{
		amountSpecifiedRemaining: amountSpecified,
		amountCalculated:         constants.Zero,
		sqrtPriceX96:             e.SqrtRatioX96,
		tick:                     e.TickCurrent,
		liquidity:                e.Liquidity,
	}

	// start swap while loop
	for state.amountSpecifiedRemaining.Cmp(constants.Zero) != 0 && state.sqrtPriceX96.Cmp(sqrtPriceLimitX96) != 0 {
		var step StepComputations
		step.sqrtPriceStartX96 = state.sqrtPriceX96

		step.tickNext, step.initialized = e.TickDataProvider.NextInitializedTickWithinOneWord(state.tick, zeroForOne, e.tickSpacing())

		if step.tickNext < uniswapv3utils.MinTick {
			step.tickNext = uniswapv3utils.MinTick
		} else if step.tickNext > uniswapv3utils.MaxTick {
			step.tickNext = uniswapv3utils.MaxTick
		}

		step.sqrtPriceNextX96, err = uniswapv3utils.GetSqrtRatioAtTick(step.tickNext)
		if err != nil {
			return nil, nil, nil, 0, err
		}
		var targetValue *big.Int
		if zeroForOne {
			if step.sqrtPriceNextX96.Cmp(sqrtPriceLimitX96) < 0 {
				targetValue = sqrtPriceLimitX96
			} else {
				targetValue = step.sqrtPriceNextX96
			}
		} else {
			if step.sqrtPriceNextX96.Cmp(sqrtPriceLimitX96) > 0 {
				targetValue = sqrtPriceLimitX96
			} else {
				targetValue = step.sqrtPriceNextX96
			}
		}

		state.sqrtPriceX96, step.amountIn, step.amountOut, step.feeAmount, err = uniswapv3utils.ComputeSwapStep(state.sqrtPriceX96, targetValue, state.liquidity, state.amountSpecifiedRemaining, e.Fee)
		if err != nil {
			return nil, nil, nil, 0, err
		}

		if exactInput {
			state.amountSpecifiedRemaining = new(big.Int).Sub(state.amountSpecifiedRemaining, new(big.Int).Add(step.amountIn, step.feeAmount))
			state.amountCalculated = new(big.Int).Sub(state.amountCalculated, step.amountOut)
		} else {
			state.amountSpecifiedRemaining = new(big.Int).Add(state.amountSpecifiedRemaining, step.amountOut)
			state.amountCalculated = new(big.Int).Add(state.amountCalculated, new(big.Int).Add(step.amountIn, step.feeAmount))
		}

		// TODO
		if state.sqrtPriceX96.Cmp(step.sqrtPriceNextX96) == 0 {
			// if the tick is initialized, run the tick transition
			if step.initialized {
				liquidityNet := e.TickDataProvider.GetTick(step.tickNext).LiquidityNet
				// if we're moving leftward, we interpret liquidityNet as the opposite sign
				// safe because liquidityNet cannot be type(int128).min
				if zeroForOne {
					liquidityNet = new(big.Int).Mul(liquidityNet, constants.NegativeOne)
				}
				state.liquidity = uniswapv3utils.AddDelta(state.liquidity, liquidityNet)
			}
			if zeroForOne {
				state.tick = step.tickNext - 1
			} else {
				state.tick = step.tickNext
			}
		} else if state.sqrtPriceX96.Cmp(step.sqrtPriceStartX96) != 0 {
			// recomphaven't moved
			state.tick, err = uniswapv3utils.GetTickAtSqrtRatio(state.sqrtPriceX96)
			if err != nil {
				return nil, nil, nil, 0, err
			}
		}
	}
	return state.amountCalculated, state.sqrtPriceX96, state.liquidity, state.tick, nil
}

func (e *EVMEdgeV3) ComputePriceImpact(amountIn *big.Int) *big.Float {
	return big.NewFloat(0)
}

func (e *EVMEdgeV3) Export() []string {
	data := struct {
		Token0       string `json:"token0"`
		Token1       string `json:"token1"`
		SqrtRatioX96 string `json:"sqrtRatioX96"`
		Liquidity    string `json:"liquidity"`
		TickCurrent  int    `json:"tickCurrent"`
		Fee          uint32 `json:"fee"`
		ZeroForOne   bool   `json:"zeroForOne"`
	}{
		Token0:       e.Token0.String(),
		Token1:       e.Token1.String(),
		SqrtRatioX96: e.SqrtRatioX96.String(),
		Liquidity:    e.Liquidity.String(),
		TickCurrent:  e.TickCurrent,
		Fee:          uint32(e.Fee),
		ZeroForOne:   e.ZeroForOne,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal edge data: %v", err)
		return []string{}
	}

	return []string{string(jsonData)}
}

func (e *EVMEdgeV3) GetWeight() *big.Float {
	return e.exchangeRate
}

func (e *EVMEdgeV3) tickSpacing() int {
	return constants.TickSpacings[e.Fee]
}
