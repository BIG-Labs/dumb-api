package grifts

import (
	"log"

	"dumb-api/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	app := actions.App()
	if app == nil {
		log.Fatal("Failed to initialize Buffalo application")
	}
	buffalo.Grifts(app.App)
}
