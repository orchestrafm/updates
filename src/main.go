package main

import (
	"github.com/orchestrafm/updates/src/database"
	"github.com/orchestrafm/updates/src/objstore"
	"github.com/orchestrafm/updates/src/routers"
	"github.com/spidernest-go/logger"
)

func main() {
	logger.Info().
		Msg("Starting up.")
	err := database.Connect()
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("MySQL Database could not be attached to.")
	}

	database.Synchronize()
	objstore.Login()
	routers.ListenAndServe()
}
