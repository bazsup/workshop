//go:build integration
// +build integration

package pocket_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/pocket"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCloudPocketsIT(t *testing.T) {
	e := echo.New()

	cfg := config.New().All()
	sql, err := sql.Open("postgres", cfg.DBConnection)
	fmt.Printf("> it test h.db %v", sql == nil)
	if err != nil {
		t.Error(err)
	}
	sql.Exec("INSERT INTO cloud_pockets (id, name, currency, balance) VALUES (12345, 'Travel Fund', 'THB', 100);")
	sql.Exec("INSERT INTO cloud_pockets (id, name, currency, balance) VALUES (67890, 'Savings', 'THB', 200);")
	hPocket := pocket.New(sql)
	e.GET("/cloud-pockets", hPocket.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	expected := `
	[
		{
			"id": 12345,
			"name": "Travel Fund",
			"currency": "THB",
			"balance": 100
		},
		{
			"id": 67890,
			"name": "Savings",
			"currency": "THB",
			"balance": 200
		}
	]
	`
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}
