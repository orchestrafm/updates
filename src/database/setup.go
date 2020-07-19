package database

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/spidernest-go/db/lib/sqlbuilder"
	"github.com/spidernest-go/db/mysql"
	"github.com/spidernest-go/logger"
	"github.com/spidernest-go/migrate"
)

var db sqlbuilder.Database

func Connect() error {
	opts := make(map[string]string)
	opts["parseTime"] = "True"
	conn := mysql.ConnectionURL{
		Database: os.Getenv("MYSQL_DB"),
		Host:     os.Getenv("MYSQL_HOST"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASS"),
		Options:  opts,
	}

	err := *new(error)
	db, err = mysql.Open(conn)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("MySQL Database Unreachable.")

		return err
	}

	return nil
}

func Synchronize() {
	versions := *new([]uint8)
	times := *new([]time.Time)
	buffers := *new([]io.Reader)
	names := *new([]string)
	box := packr.NewBox("./migrations")

	box.Walk(func(n string, f packr.File) error {
		vals := strings.Split(n, "_")

		// Assign Times
		epoch, err := strconv.ParseInt(vals[0], 10, 64)
		if err != nil {
			logger.Panic().
				Err(err).
				Msg("Integer Conversion was dropped.")
		}
		t := time.Unix(epoch, 0)

		times = append(times, t)

		// Assign Versioning
		ver, err := strconv.Atoi(vals[1])
		if err != nil {
			logger.Panic().
				Err(err).
				Msg("Embedded file, " + n + ", could not properly convert it's prefix to a version number.")
		}
		versions = append(versions, uint8(ver))

		// Assign Readers
		data, err := box.Find(n)
		if err != nil {
			logger.Panic().
				Err(err).
				Msg("Embedded file, " + n + ", does not exist, or could not be read.")
		}
		buf := bytes.NewBuffer(data)
		buffers = append(buffers, buf)

		// Assign Names
		names = append(names, n)

		return nil
	})

	if err := migrate.UpTo(versions, names, times, buffers, db); err != nil {
		logger.Panic().
			Err(err).
			Msg("Database Synchronization was unable to complete.")
	}

	logger.Info().
		Msg("Database Synchronization completed successfully.")
}
