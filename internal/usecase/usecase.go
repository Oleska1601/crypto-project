package usecase

import (
	"crypto-project/internal/usecase/repo/sqlitedb"
	"crypto-project/pkg/logger"
)

type Usecase struct {
	DB     *sqlitedb.SqliteDB
	logger *logger.Logger
}

func New(db *sqlitedb.SqliteDB, l *logger.Logger) *Usecase {
	return &Usecase{DB: db, logger: l}
}
