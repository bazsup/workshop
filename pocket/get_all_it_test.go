package pocket_test

import (
	"database/sql"
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
	if err != nil {
		t.Error(err)
	}

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
			"category": "Vacation",
			"currency": "THB",
			"balance": 100
		},
		{
			"id": 67890,
			"name": "Savings",
			"category": "Emergency Fund",
			"currency": "THB",
			"balance": 200
		}
	]
	`
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}