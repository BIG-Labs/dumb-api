package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HealthCheck(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "OK"}))
}
