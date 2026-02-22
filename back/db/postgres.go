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

	createTable(ctx)
}

func createTable(ctx context.Context) {
	query := `
	CREATE TABLE IF NOT EXISTS reservations (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		phone TEXT NOT NULL,
		date DATE NOT NULL,
		time TIME NOT NULL,
		guests INT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := Pool.Exec(ctx, query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}
