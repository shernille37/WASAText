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

	dbconn, err := sql.Open("sqlite3", filename)
	if err != nil {
		logger.WithError(err).Error("error Opening SQLite DB")
		return nil, fmt.Errorf("epening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("Database Stopping")
		_ = dbconn.Close()
	}()


	db, err := database.New(dbconn)
	if err != nil {
		logger.WithError(err).Error("error Creating AppDatabase")
		return nil, fmt.Errorf("creating AppDatabase: %w", err)
	}

	return db, nil
}