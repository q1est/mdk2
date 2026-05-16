package db

import (
	"context"
	"log"
)

func CreateTables(ctx context.Context) {

	reservationsQuery := `
	CREATE TABLE IF NOT EXISTS reservations (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		phone TEXT NOT NULL,
		date DATE NOT NULL,
		time TIME NOT NULL,
		guests INT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := Pool.Exec(ctx, reservationsQuery)
	if err != nil {
		log.Fatal("[ERROR] Failed to create reservations table:", err)
	}

	ordersQuery := `
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		phone TEXT NOT NULL,
		address TEXT NOT NULL,
		telegram TEXT NOT NULL,
		items JSONB NOT NULL,
		qty INT,
		total INT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`

	_, err = Pool.Exec(ctx, ordersQuery)
	if err != nil {
		log.Fatal("[ERROR]Failed to create orders table:", err)
	} 
	log.Println("[LOG] Tables created successfully")
}
