package model

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq" 
)

var db *sql.DB

func init() {
	connStr := "user=EvansAutoParts_admin password=admin dbname=EvansAutoParts host=localhost port=5432 sslmode=disable"

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Cannot find database. Received error:" + err.Error())
	} else {
		db = database
	}
}
