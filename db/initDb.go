package db

import (
	"log"
	"os"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Db *pgxpool.Pool

func InitDB() {
        databaseURL := os.Getenv("COMMENTO_DB_URL")
        if databaseURL == "" {
                log.Fatal("Environment variable COMMENTO_DB_URL is not set")
        }

        var err error
	Db, err = pgxpool.New(context.Background(), databaseURL)
        if err != nil {
                log.Fatalf("Unable to connect to database: %v", err)
        }
}
