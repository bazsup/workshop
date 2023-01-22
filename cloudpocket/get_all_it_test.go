//go:build integration
// +build integration

package cloudpocket_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkgo-software-engineering/workshop/cloudpocket"
	"github.com/kkgo-software-engineering/workshop/config"
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
	sql.Exec("INSERT INTO cloud_pockets (name, currency, balance) VALUES ('Travel Fund', 'THB', 100.0);")
	sql.Exec("INSERT INTO cloud_pockets (name, currency, balance) VALUES ('Savings', 'THB', 200.0);")
	hPocket := cloudpocket.New(sql)
	e.GET("/cloud-pockets", hPocket.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	fmt.Println(rec.Body)
	byteBuffer, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)

	var pockets []cloudpocket.Pocket
	err = json.Unmarshal(byteBuffer, &pockets)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Greater(t, len(pockets), 0)
}
