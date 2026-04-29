package handle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"restaurant/db"
	"restaurant/models"
)

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://q1est.github.io")
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

	
	go notifyTelegram(order)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func ReservationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://q1est.github.io")
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
}

// Создаем клиент один раз (лучше вынести в глобальную область или init)
var tgClient = &http.Client{
	Timeout: 15 * time.Second, // Общий таймаут на весь запрос
	Transport: &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second, // Конкретно на TLS
	},
}

func notifyTelegram(order models.Order) {
	token := os.Getenv("BOT_TOKEN")
	chatID := os.Getenv("CHAT_ID")
	if token == "" || chatID == "" {
		log.Println("Telegram config missing")
		return
	}

	// ... (формирование текста и jsonBody остается прежним) ...

	// Используем кастомный клиент вместо http.Post
	resp, err := tgClient.Post(
		"https://api.telegram.org/bot"+token+"/sendMessage",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		log.Println("Telegram error:", err)
		return
	}
	defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        log.Printf("Telegram API error: status %d\n", resp.StatusCode)
    }
}
