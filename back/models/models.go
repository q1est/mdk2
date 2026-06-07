package models

import (
"time"

"github.com/golang-jwt/jwt/v5"
)

// Public menu item returned to frontend
type MenuItem struct {
ID          int     `json:"id"`
Name        string  `json:"name"`
Description string  `json:"description"`
Price       float64 `json:"price"`
Category    string  `json:"category"`
ImageURL    string  `json:"image_url"`
}

// Admin/API models
type CreateMenuRequest struct {
Name        string  `json:"name"`
Description string  `json:"description"`
Price       float64 `json:"price"`
Category    string  `json:"category"`
ImageURL    string  `json:"image_url"`
Available   bool    `json:"available"`
}

type AdminMenuItem struct {
ID          int        `json:"id" db:"id"`
Name        string     `json:"name" db:"name"`
Description string     `json:"description" db:"description"`
Price       float64    `json:"price" db:"price"`
Category    string     `json:"category" db:"category"`
ImageURL    string     `json:"image_url" db:"image_url"`
Available   bool       `json:"available" db:"available"`
CreatedAt   time.Time  `json:"created_at" db:"created_at"`
UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Auth models
type LoginRequest struct {
Email    string `json:"email"`
Password string `json:"password"`
}

type User struct {
ID    string `json:"id"`
Email string `json:"email"`
Name  string `json:"name"`
}

type LoginResponse struct {
Token string `json:"token"`
User  User   `json:"user"`
}

// JWT claims for admin tokens
type JWTClaims struct {
AdminID string `json:"admin_id"`
Email   string `json:"email"`
jwt.RegisteredClaims
}
package models

type Reservation struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	Guests  int    `json:"guests"`

} 
type OrderItem struct {
	Name string `json:"name"`
	Qty  int    `json:"qty"`
}
type Order struct {
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Address  string   `json:"address"`
	Telegram string   `json:"telegram"`
	Items    []OrderItem `json:"items"`
	Total    int      `json:"total"`
}