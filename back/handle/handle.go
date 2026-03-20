package handle

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"restaurant/db"
	"restaurant/models"
)

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	ctx := context.Background()
	if r.Method != http.MethodPost {
		http.Error(w, "No Get, PLS!!!", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodPost {
		var order models.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ItemsJSON, err := json.Marshal(order.Items)
		if err != nil {
			http.Error(w, "Internal Server Error: JSON marshalling failed", http.StatusInternalServerError)
			return
		}

		_, err = db.PoolOrder.Exec(ctx, `INSERT INTO orders (name, phone, address, telegram, items, total) VALUES ($1,$2,$3,$4,$5,$6)`,
			order.Name, order.Phone, order.Address, order.Telegram, ItemsJSON, order.Total,
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
		rowsOrder, err := db.PoolOrder.Query(ctx, `SELECT name, phone, address, telegram, items, total FROM orders`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rowsOrder.Close()

		reservations := []models.Order{}
		for rowsOrder.Next() {
			var resOrde models.Order
			if err := rowsOrder.Scan(&resOrde.Name, &resOrde.Phone, &resOrde.Address, &resOrde.Telegram, &resOrde.Items, &resOrde.Total); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reservations = append(reservations, resOrde)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reservations)
	}

}
func ReservationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := context.Background()

	if r.Method != http.MethodPost {
		http.Error(w, "No Get, PLS!!!", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodPost {
		var res models.Reservation
		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := db.Pool.Exec(ctx,
			`INSERT INTO reservations (name, phone, date, time, guests) VALUES ($1,$2,$3,$4,$5)`,
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
		rows, err := db.Pool.Query(ctx, `SELECT name, phone, date::TEXT, time::TEXT, guests FROM reservations`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		reservations := []models.Reservation{}
		for rows.Next() {
			var res models.Reservation
			if err := rows.Scan(&res.Name, &res.Phone, &res.Date, &res.Time, &res.Guests); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reservations = append(reservations, res)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reservations)
	}
}
