package models

type Reservation struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	Guests  int    `json:"guests"`

}
type Order struct {
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Address  string   `json:"address"`
	Telegram string   `json:"telegram"`
	Items    []string `json:"items"`
	Total    int      `json:"total"`
}
