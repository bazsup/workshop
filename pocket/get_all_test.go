//go:build unit
// +build unit

package pocket_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/pocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCloudPockets(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		p1 := pocket.PocketModel{
			ID:       12345,
			Name:     "Travel Fund",
			Currency: "THB",
			Balance:  100.0,
		}
		p2 := pocket.PocketModel{
			ID:       67890,
			Name:     "Savings",
			Currency: "THB",
			Balance:  200.0,
		}

		db, mock, _ := sqlmock.New()
		rows := sqlmock.NewRows([]string{"id", "name", "currency", "balance"}).
			AddRow(p1.ID, p1.Name, p1.Currency, p1.Balance).
			AddRow(p2.ID, p2.Name, p2.Currency, p2.Balance)
		mock.ExpectPrepare("SELECT \\* FROM cloud_pockets").ExpectQuery().WillReturnRows(rows)
		h := pocket.New(db)

		// Assertions
		wantBody := `
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
		if assert.NoError(t, h.GetAll(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, wantBody, rec.Body.String())
		}
	})

}
