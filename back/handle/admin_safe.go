package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restaurant/db"
	"restaurant/models"
)

// AdminGetMenuSafe is a parameterized, safe version of AdminGetMenu
func AdminGetMenuSafe(w http.ResponseWriter, r *http.Request) {
	if db.Pool == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"items": []models.AdminMenuItem{}})
		return
	}

	category := r.URL.Query().Get("category")
	available := r.URL.Query().Get("available")

	query := `
		SELECT id, name, description, price, category, image_url, available, created_at, updated_at, deleted_at
		FROM menu_items
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", len(args)+1)
		args = append(args, category)
	}

	switch available {
	case "true":
		query += " AND available = true"
	case "false":
		query += " AND available = false"
	}

	query += " ORDER BY category, name"

	rows, err := db.Pool.Query(r.Context(), query, args...)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []models.AdminMenuItem{}
	for rows.Next() {
		var item models.AdminMenuItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price,
			&item.Category, &item.ImageURL, &item.Available, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt); err != nil {
			continue
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"items": items})
}
