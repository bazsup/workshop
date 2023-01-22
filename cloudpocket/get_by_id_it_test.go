//go:build integration

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
	"github.com/stretchr/testify/assert"
)

func TestGetPocketByIdIT(t *testing.T) {
	e := echo.New()
	cfg := config.New().All()
	sql, err := sql.Open("postgres", cfg.DBConnection)
	hPocket := cloudpocket.New(sql)

	if err != nil {
		t.Error(err)
	}

	stmt, err := sql.Prepare(`INSERT INTO cloud_pockets (name,currency,balance) VALUES ($1,$2,$3) RETURNING id`)
	if err != nil {
		fmt.Println("error: prepare")

	}
	row := stmt.QueryRow("shoping", "THB", 100.3)

	var id int
	fmt.Printf("id: scan%v", id)
	err = row.Scan(&id)
	if err != nil {
		fmt.Println("error: scan")
	}
	fmt.Printf("id: scan%v", id)

	e.GET("/cloud-pocket/:id", hPocket.GetPocketById)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/cloud-pocket/%d", id), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)

	var pocket cloudpocket.Pocket

	err = json.Unmarshal(body, &pocket)

	//expected := `{"name": "shoping", "currency":"THB","balance":100.3}`
	exp := cloudpocket.Pocket{
		ID:       id,
		Name:     "shoping",
		Currency: "THB",
		Balance:  100.30,
	}
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.EqualValues(t, exp.ID, pocket.ID)

	assert.EqualValues(t, exp.Name, pocket.Name)

}
