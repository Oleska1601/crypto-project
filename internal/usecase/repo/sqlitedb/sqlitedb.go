package sqlitedb

import (
	"crypto-project/pkg/logger"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteDB struct {
	database *sql.DB
}

func New(path string, logger *logger.Logger) (*SqliteDB, error) {
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	logger.Info("Open database")
	db := &SqliteDB{database: database}
	err = db.CreateTableUsers()
	logger.Info("create table users in database")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *SqliteDB) Close(logger *logger.Logger) error {
	logger.Info("Close database")
	return db.database.Close()
}
