package main

import (
	"database/sql"
	"fmt"
)

func connectToDB(dbPath string) (db *sql.DB, err error) {
	timeout := 2000 // Timeout in milliseconds
	connStr := fmt.Sprintf("file:%s?_timeout=%d", dbPath, timeout)

	db, err = sql.Open("sqlite3", connStr)
	if err != nil {
		err = fmt.Errorf("error opening database: %v", err)
		return
	}

	err = db.Ping()
	if err != nil {
		err = fmt.Errorf("error connecting to database: %v", err)
		db.Close()
		return
	}
	return
}
