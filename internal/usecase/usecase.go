package usecase

import (
	"crypto-project/internal/usecase/repo/sqlitedb"
)

type Usecase struct {
	DB *sqlitedb.SqliteDB
}

func New(db *sqlitedb.SqliteDB) *Usecase {
	return &Usecase{DB: db}
}
