package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

func DBConnect() (*sqlx.DB, error) {

	// get user + pass from env
	user := os.Getenv("CATREPORT_DB_USER")
	pass := os.Getenv("CATREPORT_DB_PASS")
	dbname := os.Getenv("CATREPORT_DB_NAME")

	// generate conn string
	connString := "postgres://" + user + ":" + pass + "@localhost/" + dbname

	// connect to the database or return error
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
