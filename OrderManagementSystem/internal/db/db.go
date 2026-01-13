package db

import (
	"database/sql"
	"log"
)

func SqlConnection() *sql.DB {
	conn := "sqlserver://username:password@localhost:1433?database=YourDB"

	db, err := sql.Open("sqlserver", conn)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	if err != nil {
		log.Fatal("connection  not estable")

	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
