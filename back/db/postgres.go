package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func ConnectPostgres() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set!")
	}

	var err error
	Pool, err = pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	if err = Pool.Ping(ctx); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	CreateTables(ctx)
}
