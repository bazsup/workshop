package pocket

import "database/sql"

type PocketModel struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

type handler struct {
	db *sql.DB
}

func New(db *sql.DB) *handler {
	return &handler{db}
}

type TransferModel struct {
	ID             int64   `json:"id"`
	PocketIDSource int     `json:"pocket_id_source"`
	PocketIDTarget int     `json:"pocket_id_target"`
	Amount         float64 `json:"amount"`
}
