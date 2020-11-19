package database

import (
	"github.com/spidernest-go/logger"
)

//HACK: I'd like to get rid of this somehow, and get back to the pureness
// of a microservice application structure
type Profile struct {
	ID                uint64   `db:"id" json:"id"`
	UUID              string   `db:"uuid" json:"uuid,omitempty"`
	Username          string   `json:"name,omitempty"`
	Groups            []string `json:"groups,omitempty"`
	Experience        uint64   `db:"experience" json:"experience"`
	Level             uint64   `db:"level" json:"level"`
	TotalScore        uint64   `db:"total_score" json:"total_score"`
	PlayCount         uint64   `db:"play_count" json:"play_count"`
	Mastery           uint8    `db:"mastery" json:"mastery"`
	PerformanceRating uint64   `db:"performance_rating" json:"performance_rating"`
}

func SelectProfileByUUID(uuid string) (error, *Profile) {
	pf := *new(Profile)
	err := db.SelectFrom("profiles").
		Where("uuid = ", uuid).
		Limit(1).
		One(&pf)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("SQL Execution had an issue when executing.")
		return err, nil
	}

	// pf.UUID = *new(string) // QUEST: Should I clear this field so it doesn't return?

	return nil, &pf
}
