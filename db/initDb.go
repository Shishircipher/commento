package db

import (
	"log"
	"os"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"fmt"
)

// InitDb initializes and returns a database connection pool.
func InitDb() (*pgxpool.Pool, error){
        databaseURL := os.Getenv("COMMENTO_DB_URL")
        if databaseURL == "" {
               log.Fatal("Environment variable COMMENTO_DB_URL is not set")
	       return nil, fmt.Errorf("environment variable COMMENTO_DB_URL is not set")
        }


	db, err := pgxpool.New(context.Background(), databaseURL)
        if err != nil {
                log.Fatalf("Unable to connect to database: %v", err)
		return nil, fmt.Errorf("unable to connect to database: %w", err)
        }
	return db, nil
}
