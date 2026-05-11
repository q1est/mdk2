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
	w.Header().Set("Access-Control-Allow-Origin", "https://q1est.github.io/mdk2/menu.htm")
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

	r.Body = http.MaxBytesReader(w, r.Body, 1024*10)

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("data received", order)
	ItemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		http.Error(w, "Internal Server Error: JSON marshalling failed", http.StatusInternalServerError)
		return
	}

	_, err = db.PoolOrder.Exec(ctx,
		`INSERT INTO orders (name, phone, address, telegram, items, total) VALUES ($1,$2,$3,$4,$5,$6)`,
		order.Name, order.Phone, order.Address, order.Telegram, ItemsJSON, order.Total,
	)

	if err != nil {
		log.Println("DB Insert error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("order сработал")
	go tg.NotifyTelegramOrder(order)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func ReservationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://q1est.github.io/mdk2/index.html")
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

	r.Body = http.MaxBytesReader(w, r.Body, 1024*10)

	var res models.Reservation
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("data received",res)
	_, err := db.Pool.Exec(ctx,
		`INSERT INTO reservations (name, phone, date, time, guests) VALUES ($1,$2,$3,$4,$5)`,
		res.Name, res.Phone, res.Date, res.Time, res.Guests,
	)
	log.Println("reservation сработал ")
	if err != nil {
		log.Println("DB Insert error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
