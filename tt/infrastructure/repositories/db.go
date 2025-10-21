package repositories

import (
	"fmt"
	"varmijo/time-tracker/tt/infrastructure/utils"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(name string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", utils.GeAppPath(fmt.Sprintf("%s.db", name)))
	if err != nil {
		return nil, err
	}

	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS records (
		id TEXT PRIMARY KEY,
		date TEXT,
		hours REAL
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS state_variables (
		key TEXT PRIMARY KEY,
		value TEXT
	)`)
	if err != nil {
		return err
	}

	return nil
}
