package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/shernille37/WASAText/service/database"
)

func startDatabase(logger logrus.FieldLogger, filename string) (database.AppDatabase, error) {

	logger.Infof("Initializing database")

	dbconn, err := sql.Open("sqlite3", filename+"?_foreign_keys=1")
	if err != nil {
		logger.WithError(err).Error("error Opening SQLite DB")
		return nil, fmt.Errorf("epening SQLite: %w", err)
	}

	db, err := database.New(dbconn)
	if err != nil {
		return nil, err
	}

	if _, err = dbconn.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, err
	}

	return db, nil
}
