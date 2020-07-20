package database

import (
	"database/sql"

	orm "github.com/spidernest-go/db"
	"github.com/spidernest-go/logger"
)

func SelectIDGreaterThan(id uint64) ([]*Patch, error) {
	updates := db.Collection("updates")
	rp := updates.Find(orm.Cond{"id": orm.Gt(id)})
	patches := *new([]*Patch)

	err := rp.All(&patches)
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
