package actions

import (
	"errors"
	"math/big"

	"dumb-api/internal/graph"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// PathRequest represents the incoming request structure
type PathRequest struct {
	TokenA  string `json:"tokenA"`
	AmountA string `json:"amountA"`
	TokenB  string `json:"tokenB"`
}

// PathResponse represents the response structure
type PathResponse struct {
	Path      []graph.Path `json:"path"`
	AmountIn  string       `json:"amountIn"`
	AmountOut string       `json:"amountOut"`
	Success   bool         `json:"success"`
	Error     string       `json:"error,omitempty"`
}

func FindBestPath(c buffalo.Context) error {
	var req PathRequest
	if err := c.Bind(&req); err != nil {
		return c.Error(400, err)
	}

	amountIn := new(big.Int)
	if _, success := amountIn.SetString(req.AmountA, 10); !success {
		return c.Error(400, errors.New("invalid amount format"))
	}

	if req.TokenA == "" || req.TokenB == "" {
		return c.Error(400, errors.New("invalid token addresses"))
	}

	updateThreshold := new(big.Float).SetFloat64(0.01)

	// Use the global graph
	g := graph.GetGlobalGraph()
	// Lock the graph
	paths := g.GetBestPaths(req.TokenA, req.TokenB, amountIn, updateThreshold)

	response := PathResponse{
		Path:     paths,
		AmountIn: req.AmountA,
		Success:  len(paths) > 0,
	}

	if len(paths) > 0 {
		response.AmountOut = paths[len(paths)-1].AmountOut.String()
	} else {
		response.Error = "no valid path found"
	}

	return c.Render(200, render.JSON(response))
}
