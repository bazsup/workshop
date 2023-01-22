//go:build integration

package cloudpocket_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kkgo-software-engineering/workshop/cloudpocket"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateItTransfer(t *testing.T) {
	e := echo.New()

	cfg := config.New().All()
	db, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Error(err)
	}

	hPocket := cloudpocket.New(db)

	e.POST("/transfers", hPocket.Transfer)

	reqBody := `{"pocket_id_source":1,"pocket_id_target":3,"amount":1.0}`
	req := httptest.NewRequest(http.MethodPost, "/transfers", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	expected := `{"id":1,"pocket_id_source":1,"pocket_id_target":3,"amount":1.0}`
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}
