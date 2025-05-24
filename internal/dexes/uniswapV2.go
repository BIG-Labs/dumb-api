package dexes

import (
	"log"
	"math/big"
	"github.com/0x7183/unifi-backend/internal/contracts"
	"github.com/0x7183/unifi-backend/internal/graph"
	"github.com/0x7183/unifi-backend/internal/graph/edges"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitUniswapV2(client *ethclient.Client, chain string, factories []string, tokens []string) (map[string]*graph.Pool, map[string]map[string]map[string]map[string]graph.Edge) {

	pools := make(map[string]*graph.Pool)
	edges := make(map[string]map[string]map[string]map[string]graph.Edge)

	for _, factoryAddr := range factories {

		factory, err := contracts.NewFactory(common.HexToAddress(factoryAddr), client)

		if err != nil {
			log.Printf("Failed to create factory client: %v", err)
			continue
		}

		for i := 0; i < len(tokens); i++ {
			for j := i + 1; j < len(tokens); j++ {
				tokenA := common.HexToAddress(tokens[i])
				tokenB := common.HexToAddress(tokens[j])

				pair, err := factory.GetPair(nil, tokenA, tokenB)

				if err != nil {
					log.Printf("Failed to get pair: %v", err)
					continue
				}

				lp, _ := contracts.NewApp(pair, client)

				token0, _ := lp.Token0(nil)

				token1, _ := lp.Token1(nil)

				reserves, _ := lp.GetReserves(nil)

				edge01 := newPoolV2EVM(token0, token1, reserves.Reserve0, reserves.Reserve1, true)
				edge10 := newPoolV2EVM(token1, token0, reserves.Reserve1, reserves.Reserve0, false)

				createUniswapV2Edges(&edges, token0, token1, pair, chain, edge01)
				createUniswapV2Edges(&edges, token1, token0, pair, chain, edge10)

				pool := &graph.Pool{
                    Token0:  token0.String(),
                    Token1:  token1.String(),
                    Pair:    pair.String(),
                    Factory: factoryAddr,
                }
                pools[pair.String()] = pool

			}
		}
	}
	return pools, edges
}

func newPoolV2EVM(token0, token1 common.Address, reserve0, reserve1 *big.Int, zeroForOne bool) graph.Edge {
	return &edges.EVMEdgeV2{
		Token0:     token0,
		Token1:     token1,
		Reserve0:   reserve0,
		Reserve1:   reserve1,
		ZeroForOne: zeroForOne,
	}
}

func createUniswapV2Edges(edges *map[string]map[string]map[string]map[string]graph.Edge, token0, token1, pair common.Address, chain string, edge graph.Edge) {

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
