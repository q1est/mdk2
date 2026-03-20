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
