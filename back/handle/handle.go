package handle

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"restaurant/db"

	"restaurant/tg"

	"restaurant/models"
)

func OrdersHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "https://q1est.github.io")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := context.Background()

	if r.Method != http.MethodPost {
		http.Error(w, "No Get, PLS!!!", http.StatusMethodNotAllowed)
		log.Print("[WARN] try method GET OrdersHAND")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*10)

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("[LOG]data orders received", order)
	ItemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		http.Error(w, "Internal Server Error: JSON marshalling failed", http.StatusInternalServerError)
		return
	}

	_, err = db.Pool.Exec(ctx,
		`INSERT INTO orders (name, phone, address, telegram, items, total) VALUES ($1,$2,$3,$4,$5,$6)`,
		order.Name, order.Phone, order.Address, order.Telegram, ItemsJSON, order.Total,
	)

	if err != nil {
		log.Println("[ERROR] DB Insert error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("[LOG] orderhand сработал")
	go tg.NotifyTelegramOrder(order)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func ReservationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://q1est.github.io")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := context.Background()

	if r.Method != http.MethodPost {
		http.Error(w, "No Get, PLS!!!", http.StatusMethodNotAllowed)
		log.Print("[WARN] try method GET ReservationsHand")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*10)

	var res models.Reservation
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("[LOG] data reservations received", res)
	_, err := db.Pool.Exec(ctx,
		`INSERT INTO reservations (name, phone, date, time, guests) VALUES ($1,$2,$3,$4,$5)`,
		res.Name, res.Phone, res.Date, res.Time, res.Guests,
	)

	if err != nil {
		log.Println("[WARN] DB Insert error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("[LOG] reservationhand сработал ")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
func GetMenu1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://q1est.github.io/mdk2/")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		log.Print("[WARN] try method POST/PUT GetMenu")
		return
	}

	items, err := db.SelectItem()
	if err != nil {
		log.Println("[ERROR] SelectItem error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}