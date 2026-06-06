package main

import (
	"log"
	"net/http"
	"os"

	"restaurant/db"
	"restaurant/handle"
	"restaurant/logs"
	"restaurant/middleware"
)

func main() {
	logs.Log()

	db.ConnectPostgres()
	defer db.Pool.Close()

	mux := http.NewServeMux()

	// Публичные API
	mux.HandleFunc("/api/orders", handle.OrdersHandler)
	mux.HandleFunc("/api/menu", handle.GetMenu)
	mux.HandleFunc("/api/reservations", handle.ReservationsHandler)

	// Админский логин
	mux.HandleFunc("/api/admin/login", handle.AdminLogin)

	// JWT Middleware
	adminAuth := middleware.JWTAuth(func(token string) (string, string, error) {
		claims, err := handle.ValidateJWT(token)
		if err != nil {
			return "", "", err
		}
		return claims.AdminID, claims.Email, nil
	})

	// Защищенные маршруты админки
	mux.Handle(
		"/api/admin/menu",
		adminAuth(http.HandlerFunc(handle.AdminGetMenu)),
	)

	// Общие middleware
	var handler http.Handler = mux

	handler = middleware.CORS(handler)
	handler = middleware.SecurityHeaders(handler)
	handler = logs.LogMiddleware(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7860"
		log.Println("[WARN] PORT not found, using default 7860")
	}

	log.Println("[LOG] Server started on :" + port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("[ERROR] Server start failed:", err)
	}
}