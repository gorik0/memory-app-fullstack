package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

type dataSources struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

func initDS() (*dataSources, error) {
	//::: POSTGRES DB
	db, err := postgresDbInit()
	if err != nil {
		return nil, err
	}
	//::: REDIS DB

	redi, err := redisDbInit()
	if err != nil {
		return nil, err
	}
	return &dataSources{DB: db, Redis: redi}, nil
}

func redisDbInit() (*redis.Client, error) {
	//	::DB_URL
	password := os.Getenv("REDIS_PASSWORD")
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	//	::DB_CONNECT

	redi := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})
	//	::DB_PING

	_, err := redi.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Error while db redis ping :::", err)

		return nil, err
	}

	return redi, nil
}

func postgresDbInit() (*sqlx.DB, error) {
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
	return db, nil
}
