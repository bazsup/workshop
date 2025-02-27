//go:build unit
// +build unit

package cloudpocket_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/cloudpocket"
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

		p1 := cloudpocket.Pocket{
			ID:       12345,
			Name:     "Travel Fund",
			Currency: "THB",
			Balance:  100.0,
		}
		p2 := cloudpocket.Pocket{
			ID:       67890,
			Name:     "Savings",
			Currency: "THB",
			Balance:  200.0,
		}

		db, mock, _ := sqlmock.New()
		rows := sqlmock.NewRows([]string{"id", "name", "parentID", "currency", "balance"}).
			AddRow(p1.ID, p1.Name, nil, p1.Currency, p1.Balance).
			AddRow(p2.ID, p2.Name, nil, p2.Currency, p2.Balance)
		mock.ExpectPrepare("SELECT \\* FROM cloud_pockets").ExpectQuery().WillReturnRows(rows)
		h := cloudpocket.New(db)

		// Assertions
		wantBody := `
		[
			{
				"id": 12345,
				"name": "Travel Fund",
				"parentID":null,
				"currency": "THB",
				"balance": 100
			},
			{
				"id": 67890,
				"name": "Savings",
				"parentID":null,
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
