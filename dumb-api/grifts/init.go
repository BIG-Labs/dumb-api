package grifts

import (
	"log"

	"github.com/0x7183/unifi-backend/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	app := actions.App()
	if app == nil {
		log.Fatal("Failed to initialize Buffalo application")
	}
	buffalo.Grifts(app.App)
}
