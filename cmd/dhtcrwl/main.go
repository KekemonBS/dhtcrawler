package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/KekemonBS/dhtcrawler/crawler"
	"github.com/KekemonBS/dhtcrawler/infrastructure/env"
	"github.com/KekemonBS/dhtcrawler/router"
	"github.com/KekemonBS/dhtcrawler/storage/postgresql"
)

func main() {
	//Init config
	cfg, err := env.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	//Init logger
	logger := log.New(os.Stdout, "log:", log.Lshortfile)

	var dbImpl crawler.DbImpl
	//Open connection
	db, err := sql.Open("postgres", cfg.PostgresURI)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Println("Postgres URI : ", cfg.PostgresURI)
	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}
	//Migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./storage/migrations",
		"postgres", driver)
	if err != nil {
		logger.Fatal(err)
	}
	m.Up()
	//Init db client
	dbImpl = postgresql.New(db)
	if err != nil {
		logger.Fatal(err)
	}

	//Init handlers
	handlers, err := crawler.New(context.Background(), logger, dbImpl, cfg.Threads)
	if err != nil {
		logger.Fatal(err)
	}

	router := router.New(handlers)
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = s.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}
