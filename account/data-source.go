package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

type dataSources struct {
	DB *sqlx.DB
}

func initDS() (*dataSources, error) {

	//	::DB_URL
	dbname := os.Getenv("POSTGRES_DBNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	sslmode := os.Getenv("POSTGRES_SSLMODE")
	user := os.Getenv("POSTGRES_USER")
	connStr := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=%s", dbname, user, password, host, port, sslmode)
	//	::DB_CONNECT
	println(connStr)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Println("Error while db connecting :::", err)

		return nil, err
	}
	//	::DB_PING

	err = db.Ping()
	if err != nil {
		log.Println("Error while db ping :::", err)

		return nil, err
	}
	//	::DB_RETURN

	return &dataSources{DB: db}, nil
}
