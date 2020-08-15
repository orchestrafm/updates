package routes

import (
	"net/http"
	"strconv"

	"github.com/orchestrafm/updates/src/database"
	"github.com/spidernest-go/logger"
	"github.com/spidernest-go/mux"
)

func listUpdates(c echo.Context) error {
	ver, err := strconv.ParseUint(c.FormValue("version"), 10, 32)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Version number was not a valid uint64")

		return c.JSON(http.StatusNotAcceptable, &struct {
			Message string
		}{
			Message: "Version was not a valid number."})
	}

	app := c.FormValue("app")
	platform := c.FormValue("platform")
	p, err := database.SelectIDGreaterThan(ver, app, platform)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Found no patches or database error")

		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, p)
}
