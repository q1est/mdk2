package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"restaurant/db"
	"restaurant/handle"
	"restaurant/logs"
	"restaurant/middleware"
)

func main() {
	logs.Log()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("[WARN] DATABASE_URL not set — skipping DB connection (development mode)")
	} else {
		db.ConnectPostgres()
		defer func() {
			if db.Pool != nil {
				db.Pool.Close()
			}
		}()
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/orders", handle.OrdersHandler)
	mux.HandleFunc("/api/menu", handle.GetMenu1)
	mux.HandleFunc("/api/reservations", handle.ReservationsHandler)

	mux.HandleFunc("/api/admin/login", handle.AdminLogin)

	adminAuth := middleware.JWTAuth(func(token string) (string, string, error) {
		claims, err := handle.ValidateJWT(token)
		if err != nil {
			return "", "", err
		}

		return claims.AdminID, claims.Email, nil
	})

	// Защищённая админка
	mux.Handle("/api/admin/menu", adminAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handle.AdminGetMenuSafe(w, r)
			return
		}
		if r.Method == http.MethodPost {
			handle.AdminCreateMenu(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})))

	mux.Handle("/api/admin/menu/", adminAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/admin/menu/")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodPut:
			handle.AdminUpdateMenu(w, r, id)
		case http.MethodDelete:
			handle.AdminDeleteMenu(w, r, id)
		case http.MethodPatch:
			handle.AdminPatchMenu(w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	mux.HandleFunc("/admin/config.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		api := os.Getenv("API_BASE_URL")
		if api == "" {
			api = "http://localhost:8080"
		}
		fmt.Fprintf(w, "window.__API_BASE_URL__ = '%s';", api)
	})

	mux.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir("./admin/"))))

	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	mux.Handle("/api/admin/upload", adminAuth(http.HandlerFunc(handle.AdminUploadImage)))

	handler := logs.LogMiddleware(middleware.CORS(mux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "7860"
	}

	log.Println("[LOG] Server started on :" + port)

	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal(err)
	}
}