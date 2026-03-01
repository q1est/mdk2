package main

import (
	"log"
	"net/http"
	"os"
	"restaurant/db"
	"restaurant/handle"
)

func main() {
	db.ConnectDBOrder()
	defer db.PoolOrder.Close()
	db.ConnectPostgres()
	defer db.Pool.Close()


	http.HandleFunc("/api/orders", handle.OrdersHandler) 
	
	http.HandleFunc("/api/reservations", handle.ReservationsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7860"
	}

	log.Println("Server started on :" + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
