package db

import "database/sql"

func initDB() (*sql.DB, error) {
	connStr := "user=biba dbname=postgres password=boba host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
