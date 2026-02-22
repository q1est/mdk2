package models 
type Reservation struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Date   string `json:"date"`
	Time   string `json:"time"`
	Guests int    `json:"guests"`
}