package pocket_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkgo-software-engineering/workshop/pocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func sqlFn() (*sql.DB, error) {
	return nil, nil
}

func TestGetAllCloudPockets(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, err := sqlFn()
		h := pocket.New(db)

		// Assertions
		wantBody := `
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
		assert.NoError(t, err)
		if assert.NoError(t, h.GetAll(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, wantBody, rec.Body.String())
		}
	})

}
