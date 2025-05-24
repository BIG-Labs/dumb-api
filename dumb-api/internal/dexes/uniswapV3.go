package dexes

import (
	"bytes"
	"context"
	"math/big"
	"time"

	"github.com/0x7183/unifi-backend/internal/contracts"
	"github.com/0x7183/unifi-backend/internal/graph"
	"github.com/0x7183/unifi-backend/internal/graph/edges"
	"github.com/0x7183/unifi-backend/internal/models"
	"github.com/0x7183/unifi-backend/internal/utils"
	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/daoleno/uniswapv3-sdk/entities"
	uniswapv3utils "github.com/daoleno/uniswapv3-sdk/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

func InitUniswapV3(client *ethclient.Client, chain string, factories []string, tokens []string, feeTiers []string) (map[string]*graph.Pool, map[string]map[string]map[string]map[string]graph.Edge) {
	pools := make(map[string]*graph.Pool)
	edges := make(map[string]map[string]map[string]map[string]graph.Edge)
	now := time.Now()

	for _, factoryAddr := range factories {
		factory, err := contracts.NewContracts(common.HexToAddress(factoryAddr), client)

		if err != nil {
			log.Printf("Failed to create factory client: %v", err)
			continue
		}

		for i := 0; i < len(tokens); i++ {
			for j := i + 1; j < len(tokens); j++ {
				tokenA := common.HexToAddress(tokens[i])
				tokenB := common.HexToAddress(tokens[j])

				for _, feeTierStr := range feeTiers {
					feeTier, ok := new(big.Int).SetString(feeTierStr, 10)

					if !ok {
						log.Printf("Failed to convert fee tier: %v", feeTierStr)
						continue
					}

					callOpts := &bind.CallOpts{
						Pending: false,
						Context: context.Background(),
					}

					pair, err := factory.GetPool(callOpts, tokenA, tokenB, feeTier)

					if err != nil {
						log.Printf("Failed to get pair: %v", err)
						continue
					}

					if pair == (common.Address{}) {
						log.Printf("No pool found for tokens %s and %s with fee tier %s", tokenA.String(), tokenB.String(), feeTierStr)
						continue
					}

					edge01, edge10, token0, token1 := CreateNewV3Pool(pair, common.HexToAddress(factoryAddr), client)

					CreateUniswapV3Edges(&edges, token0, token1, pair, chain, edge01)
					CreateUniswapV3Edges(&edges, token1, token0, pair, chain, edge10)

					pool := &graph.Pool{
						Token0:  token0.String(),
						Token1:  token1.String(),
						Pair:    pair.String(),
						Factory: factoryAddr,
					}

					pools[pair.String()] = pool
					poolState := &models.PoolState{
						ID:        uuid.Must(uuid.NewV4()),
						Token0:    token0.String(),
						Token1:    token1.String(),
						Pair:      pair.String(),
						Factory:   factoryAddr,
						ChainID:   chain,
						Status:    "active",
						CreatedAt: now,
						UpdatedAt: now,
					}

					err = models.DB.Create(poolState)
					if err != nil {
						log.Printf("Error saving pool to database: %v", err)
					}

					if edge01 != nil {
						edgeState01 := &models.EdgeState{
							ID:        uuid.Must(uuid.NewV4()),
							ChainID:   chain,
							Token0:    token0.String(),
							Token1:    token1.String(),
							PoolID:    pair.String(),
							EdgeData:  []byte(edge01.Export()[0]),
							CreatedAt: now,
							UpdatedAt: now,
						}

						err = models.DB.Create(edgeState01)
						if err != nil {
							log.Printf("Error saving edge01 to database: %v", err)
						}
					}

					if edge10 != nil {
						edgeState10 := &models.EdgeState{
							ID:        uuid.Must(uuid.NewV4()),
							ChainID:   chain,
							Token0:    token1.String(),
							Token1:    token0.String(),
							PoolID:    pair.String(),
							EdgeData:  []byte(edge10.Export()[0]),
							CreatedAt: now,
							UpdatedAt: now,
						}

						err = models.DB.Create(edgeState10)
						if err != nil {
							log.Printf("Error saving edge10 to database: %v", err)
						}
					}
				}
			}
		}
	}
	return pools, edges
}

func newPoolV3(tokenA, tokenB common.Address, fee constants.FeeAmount, sqrtRatioX96, liquidity *big.Int, tickCurrent int, ticks []entities.Tick) graph.Edge {
	if fee >= constants.FeeMax {
		return nil
	}

	token0 := tokenA
	token1 := tokenB
	isSorted := bytes.Compare(tokenA.Bytes(), tokenB.Bytes()) < 0
	if !isSorted {
		token0 = tokenB
		token1 = tokenA
	}

	provider, err := entities.NewTickListDataProvider(ticks, constants.TickSpacings[fee])
	utils.Assert(err == nil, err)

	p := &edges.EVMEdgeV3{
		Token0:           token0,
		Token1:           token1,
		SqrtRatioX96:     sqrtRatioX96,
		Liquidity:        liquidity,
		TickCurrent:      tickCurrent,
		TickDataProvider: provider,
		ZeroForOne:       isSorted,
		Fee:              constants.FeeAmount(fee),
	}

	return p
}

func CreateNewV3Pool(pairAddr, factory common.Address, client *ethclient.Client) (graph.Edge, graph.Edge, common.Address, common.Address) {
	emptyAddr := common.Address{}

	lp, err := contracts.NewPoolV3(pairAddr, client)
	utils.Assert(err == nil, err)

	token0, err := lp.Token0(nil)
	utils.Assert(err == nil, err)

	token1, err := lp.Token1(nil)
	utils.Assert(err == nil, err)

	fee, err := lp.Fee(nil)
	utils.Assert(err == nil, err)

	uintFee := constants.FeeAmount(fee.Uint64())
	tickSpacing := constants.TickSpacings[uintFee]

	if tickSpacing == 0 {
		return nil, nil, emptyAddr, emptyAddr
	}

	slot0, err := lp.Slot0(nil)
	utils.Assert(err == nil, err)

	liquidity, err := lp.Liquidity(nil)
	utils.Assert(err == nil, err)

	currentTick := int(slot0.Tick.Int64())

	var ticks []entities.Tick
	var dbTicks []models.Tick

	err = models.DB.Where("pool_address = ?", pairAddr.String()).All(&dbTicks)
	if err == nil && len(dbTicks) > 0 {
		for _, dbTick := range dbTicks {
			ticks = append(ticks, entities.Tick{
				Index:          dbTick.Index,
				LiquidityGross: dbTick.GetLiquidityGross(),
				LiquidityNet:   dbTick.GetLiquidityNet(),
			})
		}
	} else {
		minWord := utils.TickToWord(uniswapv3utils.MinTick, tickSpacing)
		maxWord := utils.TickToWord(uniswapv3utils.MaxTick, tickSpacing)

		wordPosIndices := make([]int, 0)
		results := make([]*big.Int, 0)
		i := minWord
		for i <= maxWord {
			wordPosIndices = append(wordPosIndices, i)
			var bitmapRes *big.Int
			bitmapRes, err = lp.TickBitmap(nil, int16(i))
			log.Printf("bitmapRes: %v", bitmapRes)
			if err != nil {
				log.Printf("Error fetching TickBitmap, retrying after sleep: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}
			results = append(results, bitmapRes)
			time.Sleep(250 * time.Millisecond)
			i++
		}

		tickIndices := make([]int, 0)

		for j := 0; j < len(wordPosIndices); j++ {
			ind := wordPosIndices[j]
			bitmap := results[j]

			if bitmap.Cmp(big.NewInt(0)) != 0 {
				for i := 0; i < 256; i++ {
					bit := big.NewInt(1)
					initialized := new(big.Int).And(bitmap, new(big.Int).Lsh(bit, uint(i))).Cmp(big.NewInt(0)) != 0
					if initialized {
						tickIndex := (ind*256 + i) * tickSpacing
						tickIndices = append(tickIndices, tickIndex)
					}
				}
			}
		}

		now := time.Now()
		for _, t := range tickIndices {
			tickData, err := lp.Ticks(nil, big.NewInt(int64(t)))
			log.Printf("tick: %v", tickData)
			if err != nil {
				return nil, nil, emptyAddr, emptyAddr
			}

			entityTick := entities.Tick{
				Index:          t,
				LiquidityGross: tickData.LiquidityGross,
				LiquidityNet:   tickData.LiquidityNet,
			}
			ticks = append(ticks, entityTick)

			dbTick := models.Tick{
				ID:          uuid.Must(uuid.NewV4()),
				PoolAddress: pairAddr.String(),
				Index:       t,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			dbTick.SetLiquidityGross(tickData.LiquidityGross)
			dbTick.SetLiquidityNet(tickData.LiquidityNet)

			err = models.DB.Create(&dbTick)
			if err != nil {
				log.Printf("Error saving tick to database: %v", err)
			}
		}
	}

	edge01 := newPoolV3(token0, token1, uintFee, slot0.SqrtPriceX96, liquidity, currentTick, ticks)
	edge10 := newPoolV3(token1, token0, uintFee, slot0.SqrtPriceX96, liquidity, currentTick, ticks)

	return edge01, edge10, token0, token1
}

func CreateUniswapV3Edges(edges *map[string]map[string]map[string]map[string]graph.Edge, token0, token1, pair common.Address, chain string, edge graph.Edge) {
	if _, ok := (*edges)[token0.String()]; !ok {
		(*edges)[token0.String()] = make(map[string]map[string]map[string]graph.Edge)
	}

	if _, ok := (*edges)[token0.String()][token1.String()]; !ok {
		(*edges)[token0.String()][token1.String()] = make(map[string]map[string]graph.Edge)
	}

	if _, ok := (*edges)[token0.String()][token1.String()][pair.String()]; !ok {
		(*edges)[token0.String()][token1.String()][pair.String()] = make(map[string]graph.Edge)
	}

	if _, ok := (*edges)[token0.String()][token1.String()][pair.String()][chain]; !ok {
		(*edges)[token0.String()][token1.String()][pair.String()][chain] = edge
	}
}
