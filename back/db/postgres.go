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
		log.Fatal("[ERROR] DATABASE_URL is not set!")
	} 
	log.Println("[LOG] set url")

	var err error
	Pool, err = pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal("[ERROR] Unable to connect to database:", err)
	}

	if err = Pool.Ping(ctx); err != nil {
		log.Println("[ERROR] DB not reachable:", err)
	} 
	log.Println("[LOG] DB connected successfully")

	CreateTables(ctx)
}
