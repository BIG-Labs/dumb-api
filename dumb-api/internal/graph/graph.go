package graph

import (
	"math/big"
	"sync"

	"dumb-api/config"
	"dumb-api/internal/utils"
)

// Global graph instance
var GlobalGraph *Graph

// Initialize the global graph
func InitGlobalGraph() *Graph {
	if GlobalGraph == nil {
		GlobalGraph = &Graph{
			Edges:    make(map[string]map[string]map[string]map[string]Edge),
			Pools:    make(map[string]*Pool),
			PoolFees: make(map[string]*big.Int),
		}
	}
	return GlobalGraph
}

// GetGlobalGraph returns the global graph instance
func GetGlobalGraph() *Graph {
	return GlobalGraph
}

type Graph struct {
	Mu       sync.RWMutex
	Edges    map[string]map[string]map[string]map[string]Edge
	Pools    map[string]*Pool
	PoolFees map[string]*big.Int
}

type Pool struct {
	Token0  string
	Token1  string
	Pair    string
	Factory string
}

type Path struct {
	TokenIn     string
	Pool        string
	AmountOut   *big.Int
	Chain       string
	TokenOut    string
	TokenHome   string
	TokenRemote string
}

func NewGraph() *Graph {
	return &Graph{
		Edges: make(map[string]map[string]map[string]map[string]Edge),
		Pools: make(map[string]*Pool),
	}
}

func (g *Graph) NewEdge(from, to, pool, chain string, edge Edge) {
	if g.Edges[from] == nil {
		g.Edges[from] = make(map[string]map[string]map[string]Edge)
	}
	if g.Edges[from][to] == nil {
		g.Edges[from][to] = make(map[string]map[string]Edge)
	}
	if g.Edges[from][to][pool] == nil {
		g.Edges[from][to][pool] = make(map[string]Edge)
	}
	g.Edges[from][to][pool][chain] = edge
}

func (g *Graph) NewPool(pool, token0, token1, factory string) {
	if _, exists := g.Pools[pool]; !exists {
		g.Pools[pool] = &Pool{
			Token0:  token0,
			Token1:  token1,
			Pair:    pool,
			Factory: factory,
		}
	}
}

func (g *Graph) HasPool(pool string) bool {
	_, exists := g.Pools[pool]
	return exists
}

func (g *Graph) GetPool(pool string) *Pool {
	return g.Pools[pool]
}

func (g *Graph) GetPoolFee(pool string) *big.Int {
	if fee, exists := g.PoolFees[pool]; exists {
		return fee
	}
	return config.DefaultBuilderFee
}

func (g *Graph) GetPools(token0, token1 string) map[string]map[string]Edge {
	if _, exists := g.Edges[token0]; exists {
		if p, exists := g.Edges[token0][token1]; exists {
			return p
		}
	}
	return nil
}

func (g *Graph) GetEdge(token0, token1, pool, chain string) Edge {
	if _, exists := g.Edges[token0]; exists {
		if _, exists := g.Edges[token0][token1]; exists {
			if edge, exists := g.Edges[token0][token1][pool]; exists {
				if edge, exists := edge[chain]; exists {
					return edge
				}
			}
		}
	}
	return nil
}

func (g *Graph) DeleteEdge(token0, token1, pool string) {
	if _, exists := g.Edges[token0]; exists {
		if _, exists := g.Edges[token0][token1]; exists {
			if _, exists := g.Edges[token0][token1][pool]; exists {
				delete(g.Edges[token0][token1], pool)
				delete(g.Edges[token1][token0], pool)
			}
		}
	}
}

func (g *Graph) DeletePool(pool string) {
	pair := g.Pools[pool]
	delete(g.Pools, pool)
	delete(g.Edges[pair.Token0][pair.Token1], pool)
	delete(g.Edges[pair.Token1][pair.Token0], pool)
}

func (g *Graph) GetBestPaths(start, target string, amountIn *big.Int, updateThreshold *big.Float) []Path {
	g.Mu.RLock()
	defer g.Mu.RUnlock()

	paths := make(map[string]Path)
	visitedPools := make(map[string]bool)

	queue := utils.NewQueue[string]()
	queue.Enqueue(start)

	paths[start] = Path{
		TokenIn:   start,
		TokenOut:  "",
		AmountOut: amountIn,
		Chain:     "avalanche",
	}

	for !queue.IsEmpty() {
		currentAsset := queue.Dequeue()

		currentBestAmount := paths[currentAsset].AmountOut

		for targetAsset, pools := range g.Edges[currentAsset] {
			var (
				bestPoolAmount *big.Int
				bestPool       string
				bestChain      string
			)

			for pool, edges := range pools {
				if _, ok := visitedPools[pool]; ok {
					continue
				}

				for chain, edge := range edges {
					highSa := edge.ComputeExactAmountOut(currentBestAmount)

					if highSa == nil || highSa.Sign() <= 0 {
						continue
					}

					if bestPoolAmount == nil || highSa.Cmp(bestPoolAmount) > 0 {
						bestPoolAmount = highSa
						bestPool = pool
						bestChain = chain
					}
				}
			}

			if bestPoolAmount != nil {
				visitedPools[bestPool] = true
			}

			if bestPoolAmount == nil {
				continue
			}

			queue.Enqueue(targetAsset)

			if existingPath, exists := paths[targetAsset]; !exists || bestPoolAmount.Cmp(existingPath.AmountOut) > 0 {
				tokenHome := g.GetTokenHome(targetAsset, bestChain)
				tokenRemote := g.GetTokenRemote(targetAsset, bestChain)

				paths[targetAsset] = Path{
					TokenIn:     currentAsset,
					Pool:        bestPool,
					AmountOut:   bestPoolAmount,
					Chain:       bestChain,
					TokenOut:    targetAsset,
					TokenHome:   tokenHome,
					TokenRemote: tokenRemote,
				}
			}
		}
	}

	return buildPath(paths, start, target)
}

func buildPath(paths map[string]Path, start, target string) []Path {
	var path []Path
	for target != start {
		path = append([]Path{paths[target]}, path...)
		target = paths[target].TokenIn
		if target == "" {
			break
		}
	}
	return path
}

func (g *Graph) GetTokenHome(token, chain string) string {
	home, _ := utils.GetTokenMappings(token, chain)
	return home
}

func (g *Graph) GetTokenRemote(token, chain string) string {
	_, remote := utils.GetTokenMappings(token, chain)
	return remote
}
