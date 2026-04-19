package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var PoolOrder *pgxpool.Pool

func ConnectDBOrder() {
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL_ORDER")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set!")
	}

	var err error
	PoolOrder, err = pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	if err = PoolOrder.Ping(ctx); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	createTableOrder(ctx)
}

func createTableOrder(ctx context.Context) {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		phone TEXT NOT NULL,
		address TEXT NOT NULL,
		telegram TEXT NOT NULL,
		items JSONB NOT NULL,
		qty INT NOT NULL,
		total INT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()

);`
	_, err := PoolOrder.Exec(ctx, query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

}
