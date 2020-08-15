package routes

import (
	"net/http"

	"github.com/orchestrafm/updates/src/database"
	"github.com/orchestrafm/updates/src/objstore"
	"github.com/spidernest-go/logger"
	"github.com/spidernest-go/mux"
)

func pushUpdate(c echo.Context) error {
	if authorized := HasRole(c, "create-update"); authorized != true {
		logger.Info().
			Msg("user intent to push a patch, but was unauthorized.")

		return c.JSON(http.StatusUnauthorized, &struct {
			Message string
		}{
			Message: ErrPermissions.Error()})
	}

	// Data Binding
	p := new(database.Patch)
	p.Application = c.FormValue("app")
	p.Name = c.FormValue("name")
	p.Platform = c.FormValue("platform")

	f, err := c.FormFile("patch")
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Multipart Form data file is invalid or missing.")

		return c.JSON(http.StatusNotAcceptable, &struct {
			Message string
		}{
			Message: "Patch file was invalid or missing."})
	}

	src, err := f.Open()
	defer src.Close()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Multipart Form data file is malformed.")

		return c.JSON(http.StatusNotAcceptable, &struct {
			Message string
		}{
			Message: "Patch file is malformed."})
	}

	url, err := objstore.Upload(src, f.Filename)
	p.URL = url
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Object Storage rejected putting the object.")

		return c.JSON(http.StatusInternalServerError, &struct {
			Message string
		}{
			Message: "File could not be commited to disk."})
	}

	err = p.New()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Patch could not get journaled to the database and is now inconsistent with Object Storage permenence.")

		return c.JSON(http.StatusNotAcceptable, &struct {
			Message string
		}{
			Message: "Patch was not written to the database."})
	}

	return c.JSON(http.StatusOK, &struct {
		Message string
	}{
		Message: "Patch uploaded successfully."})
}
