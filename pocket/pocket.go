package pocket

import "database/sql"

type PocketModel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Currency string `json:"currency"`
	Balance  int64  `json:"balance"`
}

type handler struct {
	db *sql.DB
}

func New(db *sql.DB) *handler {
	return &handler{db}
}
