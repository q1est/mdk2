package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"restaurant/models"
	"strings"
	"net/http"
	"os"
)

func NotifyTelegramOrder(order models.Order) {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Println("[WARN] BOT_TOKEN not set")
		return
	}

	chatID := os.Getenv("CHAT_ID")
	if chatID == "" {
		log.Println("[WARN]CHAT_ID not set")
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
	var sb strings.Builder
	for _, item := range order.Items {
		fmt.Fprintf(&sb, "• %s x%d\n", item.Name, item.Qty)
	}

	text += fmt.Sprintf("\n💰 Итого: %d₽", order.Total)

	body := map[string]any{
		"chat_id": chatID,
		"text":    sb.String(),
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(
		"https://api.telegram.org/bot"+token+"/sendMessage",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		log.Println("[WARN]Telegram error:", err)
		return
	} 
	log.Println("[LOG]tg bot work")
	defer resp.Body.Close()
}
