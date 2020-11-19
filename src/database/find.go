package database

import (
	"database/sql"

	orm "github.com/spidernest-go/db"
	"github.com/spidernest-go/logger"
)

func SelectIDGreaterThan(id uint64, app, platform string) ([]*Patch, error) {
	updates := db.Collection("updates")
	rp := updates.Find()
	patches := *new([]*Patch)

	err := rp.Where("app = ", app).
		And("platform = ", platform).
		And(orm.Cond{"id": orm.Gt(id)}).
		All(&patches)
	if err != nil && err != sql.ErrNoRows {
		logger.Error().
			Err(err).
			Msg("Bad parameters or database error.")

		return nil, err
	}

	if err == sql.ErrNoRows {
		logger.Info().
			Err(err).
			Msg("No patches found after specified ID.")

		return nil, err
	}

	return patches, nil
}
