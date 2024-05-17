package db

import (
	"database/sql"
	"fmt"
	"log"
)

func ConnectDB() (*sql.DB, error) {
	// BD connection
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "123"
		dbname   = "Quizlex"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
