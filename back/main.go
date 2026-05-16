package main

import (
	"log"
	"net/http"
	"os"
	"restaurant/db"
	"restaurant/handle"
	"restaurant/logs"
)

func main() { 
	logs.Log()
	db.ConnectPostgres()
	defer db.Pool.Close()
	
mux := http.NewServeMux()
	mux.HandleFunc("/api/orders", handle.OrdersHandler) 
	
	mux.HandleFunc("/api/reservations", handle.ReservationsHandler) 
	logmux := logs.LogMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7860" 
		log.Println("[WARN], not get port, default port 7860")
	}

	log.Println("Server started on :" + port, )
	err := http.ListenAndServe(":"+port, logmux)
	if err != nil {
		log.Fatal("[ERROR] not start",err)
	} 
	
//	log.Fatal(http.ListenAndServe(":"+port, nil, )) 
	//http.ListenAndServe(""+port,nil) 

}
