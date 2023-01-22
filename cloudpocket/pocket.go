package cloudpocket

import (
	"database/sql"
	"math"
)

type Pocket struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	ParentID *int    `json:"parentID"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

type handler struct {
	db *sql.DB
}

func New(db *sql.DB) *handler {
	return &handler{db}
}

type Transfer struct {
	ID             int     `json:"id"`
	PocketIDSource int     `json:"pocket_id_source"`
	PocketIDTarget int     `json:"pocket_id_target"`
	Amount         float64 `json:"amount"`
}

func (t *Transfer) IsValidAmount() bool {
	return t.Amount >= 0.01 && t.isValidDecimal()
}

func (t *Transfer) isValidDecimal() bool {
	return round(t.Amount) == t.Amount
}

func round(num float64) float64 {
	return math.Round(num*100) / 100.0
}
