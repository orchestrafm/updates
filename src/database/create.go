package database

import (
	"github.com/spidernest-go/logger"
)

func (p *Patch) New() error {
	r, err := db.InsertInto("updates").
		Values(p).
		Exec()

	if err != nil {
		logger.Error().
			Err(err).
			Msg("Patch could not be inserted into the table.")
	}

	id, err := r.LastInsertId()
	if err == nil {
		p.ID = uint64(id)
	}
	return err
}
