package cloudpocket

import "database/sql"

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
