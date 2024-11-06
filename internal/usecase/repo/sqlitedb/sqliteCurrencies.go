package sqlitedb

import (
	"crypto-project/internal/entity"
	"fmt"
)

func (db *SqliteDB) CreateTableResponses() error {
	createTableResponses := `CREATE TABLE IF NOT EXISTS responses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INT,
    BTC REAL,
    ETH REAL,
    LTC REAL,
    USDT REAL,
    FOREIGN KEY (user_id) REFERENCES users(id)                                 
)`
	_, err := db.database.Exec(createTableResponses)
	return err
}

func (db *SqliteDB) InsertTableResponses(user_id int, response *entity.Response) error {
	insertTableUsers := `INSERT INTO responses (id, user_id, BTC, ETH, LTC, USDT) VALUES (NULL, ?, ?, ?, ?, ?)`
	_, err := db.database.Exec(insertTableUsers, user_id, response.BTC, response.ETH, response.LTC, response.USDT)
	if err != nil {
		return fmt.Errorf("database Exec error: %v", err)
	}
	return nil
}

func (db *SqliteDB) GetTableResponses(user_id int) ([]entity.Response, error) {
	userExists := `SELECT BTC, ETH, LTC, USDT FROM responses WHERE user_id =?`
	rows, err := db.database.Query(userExists, user_id)
	if err != nil {
		return nil, fmt.Errorf("database Query error: %v", err)
	}
	var databaseResponses []entity.Response
	defer rows.Close()
	for rows.Next() {
		var response entity.Response
		err = rows.Scan(&response.BTC, &response.ETH, &response.LTC, &response.USDT)
		if err != nil {
			return nil, fmt.Errorf("rows Scan error: %v", err)
		}
		databaseResponses = append(databaseResponses, response)
	}
	return databaseResponses, nil
}
