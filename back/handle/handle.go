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

	// ✅ уведомление в Telegram
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

func notifyTelegram(order models.Order) {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Println("BOT_TOKEN not set")
		return
	}

	chatID := os.Getenv("CHAT_ID")
	if chatID == "" {
		log.Println("CHAT_ID not set")
		return
	}

	text := fmt.Sprintf(
		"🆕 Новый заказ\n\n"+
			"👤 %s\n"+
			"📞 %s\n"+
			"🏠 %s\n\n"+
			"🍽 Блюда:\n",
		order.Name,
		order.Phone,
		order.Address,
	)

	for _, item := range order.Items {
		text += fmt.Sprintf("• %s x%d\n", item.Name, item.Qty)
	}

	text += fmt.Sprintf("\n💰 Итого: %d₽", order.Total)

	body := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(
		"https://api.telegram.org/bot"+token+"/sendMessage",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		log.Println("Telegram error:", err)
		return
	}
	defer resp.Body.Close()
}
