package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Reservation struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Date   string `json:"date"`
	Time   string `json:"time"`
	Guests int    `json:"guests"`
}

var db *pgxpool.Pool

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")

	var err error
	db, err = pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping(ctx)
	if err != nil {
		log.Fatal("DB not reachable:", err)
	}

	createTable(ctx)

	http.HandleFunc("/api/reservations", reservationsHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
}

func reservationsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")


	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := context.Background()

	if r.Method == http.MethodPost {
		var res Reservation
		err := json.NewDecoder(r.Body).Decode(&res)
		if err != nil {
			log.Println("Decode error:", err) // Логируем для отладки
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec(ctx,
			`INSERT INTO reservations (name, phone, date, time, guests)
			 VALUES ($1,$2,$3,$4,$5)`,
			res.Name, res.Phone, res.Date, res.Time, res.Guests,
		)

		if err != nil {
			log.Println("DB Insert error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	if r.Method == http.MethodGet {
		rows, err := db.Query(ctx,
			 `SELECT name, phone, date::TEXT, time::TEXT, guests FROM reservations`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Чтобы не возвращать null, если база пуста, инициализируем слайс
		reservations := []Reservation{} 

		for rows.Next() {
			var r Reservation
			err := rows.Scan(&r.Name, &r.Phone, &r.Date, &r.Time, &r.Guests)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reservations = append(reservations, r)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reservations)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	if (*w).Header().Get("Access-Control-Request-Method") != "" {
		(*w).WriteHeader(http.StatusOK)
	}
}
