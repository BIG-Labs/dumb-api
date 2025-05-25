package main

import (
	"log"

	"github.com/0x7183/unifi-backend/actions"
	"github.com/0x7183/unifi-backend/config"
	"github.com/0x7183/unifi-backend/internal/dexes"
	"github.com/0x7183/unifi-backend/internal/graph"
	"github.com/0x7183/unifi-backend/internal/graph/edges"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func initializeGraph() {
	globalGraph := graph.InitGlobalGraph()

	client, err := ethclient.Dial(config.AVALANCHE_RPC_URL)
	if err != nil {
		log.Fatalf("Failed to connect to %s network: %v", "AVALANCHE", err)
	}

	pools, v3Edges := dexes.InitUniswapV3(
		client,
		"avalanche",
		config.EVMConfig["AVALANCHE"].UniswapV3.Factories,
		config.TokensByChain["AVALANCHE"],
		config.FeeTiers,
	)

	globalGraph.Mu.Lock()

	for addr, pool := range pools {
		globalGraph.NewPool(addr, pool.Token0, pool.Token1, pool.Factory)
	}

	globalGraph.Edges = v3Edges

	createStaticPool(globalGraph)

	globalGraph.Mu.Unlock()

}

func createStaticPool(globalGraph *graph.Graph) {
	usdcAvalanche := config.EVMConfig["AVALANCHE"].Tokens[0].Address
	usdcCoq := config.EVMConfig["COQNET"].Tokens[0].Address

	bridgeEdge01 := edges.BridgeEdge{
		Token0: common.HexToAddress(usdcAvalanche),
		Token1: common.HexToAddress(usdcCoq),
	}

	bridgeEdge10 := edges.BridgeEdge{
		Token0: common.HexToAddress(usdcCoq),
		Token1: common.HexToAddress(usdcAvalanche),
	}

	globalGraph.NewEdge(usdcAvalanche, usdcCoq, "pool", "coqnet", &bridgeEdge01)
	globalGraph.NewEdge(usdcCoq, usdcAvalanche, "pool1", "avalanche", &bridgeEdge10)

}

// main is the starting point for your Buffalo application.
// You can feel free and add to this `main` method, change
// what it does, etc...
// All we ask is that, at some point, you make sure to
// call `app.Serve()`, unless you don't want to start your
// application that is. :)
func main() {
	initializeGraph()

	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}

/*
# Notes about `main.go`

## SSL Support

We recommend placing your application behind a proxy, such as
Apache or Nginx and letting them do the SSL heavy lifting
for you. https://gobuffalo.io/en/docs/proxy

## Buffalo Build

When `buffalo build` is run to compile your binary, this `main`
function will be at the heart of that binary. It is expected
that your `main` function will start your application using
the `app.Serve()` method.

*/
