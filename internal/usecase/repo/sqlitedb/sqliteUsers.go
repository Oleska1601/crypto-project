package sqlitedb

import (
	"crypto-project/internal/entity"
	"fmt"
)

func (db *SqliteDB) CreateTableUsers() error {
	createTableUsers := `CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login TEXT,
    password_hash TEXT,
    salt TEXT,
    secret TEXT
)`
	_, err := db.database.Exec(createTableUsers)
	return err
}

func (db *SqliteDB) InsertTableUsers(user *entity.User) error {
	insertTableUsers := `INSERT INTO users (id, login, password_hash, salt, secret) VALUES (NULL, ?, ?, ?, ?)`
	_, err := db.database.Exec(insertTableUsers, user.Login, user.PasswordHash, user.Salt, user.Secret)
	if err != nil {
		return fmt.Errorf("database Exec error: %v", err)
	}
	return nil
}

//если пользователь существует, будут возвращены его параметры, в противном случае пустая структура

func (db *SqliteDB) GetTableUsers(user *entity.User) (entity.User, error) {
	userExists := `SELECT id, login, password_hash, salt, secret FROM users WHERE login =?`
	row, err := db.database.Query(userExists, user.Login)
	databaseUser := entity.User{}
	if err != nil {
		return databaseUser, fmt.Errorf("database Query error: %v", err)
	}
	defer row.Close()
	if row.Next() {
		err = row.Scan(&databaseUser.ID, &databaseUser.Login, &databaseUser.PasswordHash, &databaseUser.Salt, &databaseUser.Secret)
		if err != nil {
			return databaseUser, fmt.Errorf("row Scan error: %v", err)
		}
	}

	return databaseUser, nil

}
