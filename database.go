package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

func DBConnect() (*sql.DB, error) {

	// get user + pass from env
	user := os.Getenv("CATREPORT_DB_USER")
	pass := os.Getenv("CATREPORT_DB_PASS")
	dbname := os.Getenv("CATREPORT_DB_NAME")

	// generate conn string
	connString := "postgres://" + user + ":" + pass + "@localhost/" + dbname

	// connect to the database or return error
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
