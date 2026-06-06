package db

import (
	"context"
	"log"
	"os"
	"restaurant/models"

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
func SelectItem() ([]models.MenuItem, error) {
	rows, err := Pool.Query(
		context.Background(),
		`SELECT
    id,
    name,
    description,
    price,
    category,
    image_url
	FROM menu_items
	WHERE available = true
	`,
	)
	if err != nil {
		log.Println("[ERROR] Query error:", err)
		return nil, err
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Category, &item.ImageURL)
		if err != nil {
			log.Println("[ERROR] Scan error:", err)
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		log.Println("[ERROR] Rows error:", err)
		return nil, err
	}

	return items, nil
}
