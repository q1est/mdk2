package models

type Reservation struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Date   string `json:"date"`
	Time   string `json:"time"`
	Guests int    `json:"guests"`
}
type OrderItem struct {
	Name string `json:"name"`
	Qty  int    `json:"qty"`
}
type Order struct {
	Name     string      `json:"name"`
	Phone    string      `json:"phone"`
	Address  string      `json:"address"`
	Telegram string      `json:"telegram"`
	Items    []OrderItem `json:"items"`
	Total    int         `json:"total"`
}

type MenuItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
}
