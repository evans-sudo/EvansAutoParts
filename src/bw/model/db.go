package model

import (
	"database/sql"
	"log"
)

var db *sql.DB


func init() {
	database, err := sql.Open("", "")
	if err != nil {
		log.Fatal("Cannot find database. Received error:" + err.Error())
	} else {
		db = database
	}
}