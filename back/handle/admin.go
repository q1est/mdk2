package handle

import (
"crypto/sha256"
"encoding/hex"
"encoding/json"
"net/http"
"os"
"restaurant/db"
"restaurant/models"
"strconv"
"time"

"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AdminLogin(w http.ResponseWriter, r *http.Request) {
var req models.LoginRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request body", http.StatusBadRequest)
return
}

if req.Email == "" || req.Password == "" {
http.Error(w, "Email and password required", http.StatusBadRequest)
return
}

adminEmail := os.Getenv("ADMIN_EMAIL")
adminPassword := os.Getenv("ADMIN_PASSWORD")

if req.Email != adminEmail || req.Password != adminPassword {
http.Error(w, "Invalid credentials", http.StatusUnauthorized)
return
}

token, err := GenerateJWT("1", adminEmail, 24*time.Hour)
if err != nil {
http.Error(w, "Token generation error", http.StatusInternalServerError)
return
}

resp := models.LoginResponse{
Token: token,
User: models.User{
ID:    "1",
Email: adminEmail,
Name:  "Administrator",
},
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(resp)
}

func GenerateJWT(adminID, email string, expiry time.Duration) (string, error) {
claims := models.JWTClaims{
AdminID: adminID,
Email:   email,
RegisteredClaims: jwt.RegisteredClaims{
ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
IssuedAt:  jwt.NewNumericDate(time.Now()),
Issuer:    "mdk2-admin",
},
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*models.JWTClaims, error) {
claims := &models.JWTClaims{}
token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
return jwtSecret, nil
})

if err != nil || !token.Valid {
return nil, err
}

return claims, nil
}

func HashPassword(password string) string {
hash := sha256.Sum256([]byte(password))
return hex.EncodeToString(hash[:])
}

func CheckPassword(hash, password string) bool {
return hash == HashPassword(password)
}

func AdminCreateMenu(w http.ResponseWriter, r *http.Request) {
var req models.CreateMenuRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request body", http.StatusBadRequest)
return
}

if req.Name == "" || req.Price <= 0 || req.Category == "" {
http.Error(w, "Name, price, and category are required", http.StatusBadRequest)
return
}

query := `
INSERT INTO menu_items (name, description, price, category, image_url, available, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id, created_at, updated_at
`

var newID int
var createdAt, updatedAt time.Time
err := db.Pool.QueryRow(r.Context(), query, req.Name, req.Description, req.Price, req.Category, req.ImageURL, req.Available).Scan(&newID, &createdAt, &updatedAt)
if err != nil {
http.Error(w, "Database error", http.StatusInternalServerError)
return
}

item := models.AdminMenuItem{
ID:          newID,
Name:        req.Name,
Description: req.Description,
Price:       req.Price,
Category:    req.Category,
ImageURL:    req.ImageURL,
Available:   req.Available,
CreatedAt:   createdAt,
UpdatedAt:   updatedAt,
}

w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(map[string]interface{}{"item": item})
}

func AdminUpdateMenu(w http.ResponseWriter, r *http.Request, id string) {
var req models.CreateMenuRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request body", http.StatusBadRequest)
return
}

itemID, err := strconv.Atoi(id)
if err != nil {
http.Error(w, "Invalid ID", http.StatusBadRequest)
return
}

query := `
UPDATE menu_items
SET name = $1, description = $2, price = $3, category = $4, image_url = $5, available = $6, updated_at = NOW()
WHERE id = $7 AND deleted_at IS NULL
`

result, err := db.Pool.Exec(r.Context(), query, req.Name, req.Description, req.Price, req.Category, req.ImageURL, req.Available, itemID)
if err != nil {
http.Error(w, "Database error", http.StatusInternalServerError)
return
}

if result.RowsAffected() == 0 {
http.Error(w, "Menu item not found", http.StatusNotFound)
return
}

item := models.AdminMenuItem{
ID:          itemID,
Name:        req.Name,
Description: req.Description,
Price:       req.Price,
Category:    req.Category,
ImageURL:    req.ImageURL,
Available:   req.Available,
UpdatedAt:   time.Now(),
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]interface{}{"item": item})
}

func AdminDeleteMenu(w http.ResponseWriter, r *http.Request, id string) {
itemID, err := strconv.Atoi(id)
if err != nil {
http.Error(w, "Invalid ID", http.StatusBadRequest)
return
}

query := `UPDATE menu_items SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

result, err := db.Pool.Exec(r.Context(), query, itemID)
if err != nil {
http.Error(w, "Database error", http.StatusInternalServerError)
return
}

if result.RowsAffected() == 0 {
http.Error(w, "Menu item not found", http.StatusNotFound)
return
}

w.WriteHeader(http.StatusNoContent)
}

func AdminPatchMenu(w http.ResponseWriter, r *http.Request, id string) {
var req map[string]interface{}
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request body", http.StatusBadRequest)
return
}

itemID, err := strconv.Atoi(id)
if err != nil {
http.Error(w, "Invalid ID", http.StatusBadRequest)
return
}

if available, ok := req["available"].(bool); ok {
query := `UPDATE menu_items SET available = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`
_, err := db.Pool.Exec(r.Context(), query, available, itemID)
if err != nil {
http.Error(w, "Database error", http.StatusInternalServerError)
return
}

item := models.AdminMenuItem{
ID:        itemID,
Available: available,
UpdatedAt: time.Now(),
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]interface{}{"item": item})
} else {
http.Error(w, "Invalid request", http.StatusBadRequest)
}
}
