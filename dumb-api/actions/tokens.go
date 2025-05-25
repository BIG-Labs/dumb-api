package actions

import (
	"net/http"

	"dumb-api/internal/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop/v6"
)

func GetTokens(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	tokens := &models.Tokens{}

	if err := tx.All(tokens); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, render.JSON(tokens))
}
