package actions

import (
	"log"

	"github.com/0x7183/unifi-backend/internal/graph"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v6"
)

type AppContext struct {
	*buffalo.App
	graph *graph.Graph
}

var appInstance *AppContext

func App() *AppContext {
	if appInstance == nil {
		buffApp := buffalo.New(buffalo.Options{
			Env:         envy.Get("GO_ENV", "development"),
			SessionName: "_unifi_backend_session",
		})

		db, err := pop.Connect(envy.Get("GO_ENV", "development"))
		if err != nil {
			log.Fatalf("Could not connect to database: %v", err)
		}

		buffApp.Use(func(next buffalo.Handler) buffalo.Handler {
			return func(c buffalo.Context) error {
				c.Set("tx", db)
				return next(c)
			}
		})

		g := graph.NewGraph()
		appInstance = &AppContext{
			App:   buffApp,
			graph: g,
		}

		appInstance.GET("/api/v1/prices", GetPriceData)
		appInstance.POST("/api/v1/path", FindBestPath)
		appInstance.GET("/api/v1/tokens", GetTokens)
	}
	return appInstance
}
