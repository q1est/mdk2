package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const AdminIDKey contextKey = "admin_id"
const AdminEmailKey contextKey = "admin_email"

func JWTAuth(validateFunc func(string) (string, string, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			adminID, email, err := validateFunc(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AdminIDKey, adminID)
			ctx = context.WithValue(ctx, AdminEmailKey, email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAdminID(r *http.Request) string {
	id := r.Context().Value(AdminIDKey)
	if id == nil {
		return ""
	}
	return id.(string)
}

func GetAdminEmail(r *http.Request) string {
	email := r.Context().Value(AdminEmailKey)
	if email == nil {
		return ""
	}
	return email.(string)
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8080",
			"https://mdk2.ru",
			"https://q1est.github.io/mdk2", 
			"https://q1est.github.io",
			"https://q1est.github.io/mdk2",
			"https://q1est.github.io/mdk2/", 
			"https://qwefsdfsdsg-mdk.hf.space",
		}

		allowed := false
		for _, o := range allowedOrigins {
			if origin == o {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' https://cdn.jsdelivr.net; style-src 'self' https://cdn.jsdelivr.net 'unsafe-inline'")

		next.ServeHTTP(w, r)
	})
}
