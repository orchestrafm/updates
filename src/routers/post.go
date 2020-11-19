package routers

import (
	"bytes"
	"hash/crc32"
	"net/http"

	"github.com/orchestrafm/updates/src/database"
	"github.com/orchestrafm/updates/src/objstore"
	"github.com/spidernest-go/logger"
	"github.com/spidernest-go/mux"
)

func pushUpdate(c echo.Context) error {
	if authorized := HasRole(c, "create-patch"); authorized != true {
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

	fp, err := c.FormFile("patch")
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Multipart Form data file is invalid or missing.")

		return c.JSON(http.StatusNotAcceptable, &struct {
			Message string
		}{
			Message: "Patch file was invalid or missing."})
	}

	fsi, err := c.FormFile("signature")
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Multipart Form data file is invalid or missing.")

		return c.JSON(http.StatusNotAcceptable, &struct {
			Message string
		}{
			Message: "Signature file was invalid or missing."})
	}

	// Find Issuer
	user := SelfAuthCheck(c)
	err, pf := database.SelectProfileByUUID(user.Subject)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Issuing user's profile could not be found.")

		return c.JSON(http.StatusInternalServerError, &struct {
			Message string
		}{
			Message: "Profile for issuing user went missing."})
	}
	p.Issuer = pf.ID

	// Open Multipart Files
	pmpf, err := fp.Open()
	defer pmpf.Close()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Multipart Form data file is malformed.")

		return c.JSON(http.StatusNotAcceptable, &struct {
			Message string
		}{
			Message: "Patch file is malformed."})
	}

	smpf, err := fsi.Open()
	defer smpf.Close()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Multipart Form data file is malformed.")

		return c.JSON(http.StatusNotAcceptable, &struct {
			Message string
		}{
			Message: "Patch file is malformed."})
	}

	// Check Hashes
	crc32c := crc32.MakeTable(crc32.Castagnoli)

	buf := bytes.NewBuffer([]byte{})
	buf.ReadFrom(pmpf)
	p.Hash = crc32.Checksum(buf.Bytes(), crc32c)

	buf = bytes.NewBuffer([]byte{})
	buf.ReadFrom(smpf)
	p.SignatureHash = crc32.Checksum(buf.Bytes(), crc32c)

	// Upload Files
	url, err := objstore.Upload(pmpf, fp.Filename)
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

	url, err = objstore.Upload(smpf, fsi.Filename)
	p.Signature = url
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Object Storage rejected putting the object.")

		return c.JSON(http.StatusInternalServerError, &struct {
			Message string
		}{
			Message: "File could not be commited to disk."})
	}

	// Submit Patch To Database
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
