//go:build integration

package cloudpocket_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	row1 := db.QueryRow("INSERT INTO cloud_pockets (name, currency, balance) VALUES ('Travel Fund', 'THB', 0.1) RETURNING id;")
	row2 := db.QueryRow("INSERT INTO cloud_pockets (name, currency, balance) VALUES ('Savings', 'THB', 0.2) RETURNING id;")

	var id1, id2 int
	row1.Scan(&id1)
	row2.Scan(&id2)

	hPocket := cloudpocket.New(db)

	e.POST("/transfers", hPocket.Transfer)

	reqBody := fmt.Sprintf(`{"pocket_id_source":%d,"pocket_id_target":%d,"amount":0.1}`, id1, id2)
	req := httptest.NewRequest(http.MethodPost, "/transfers", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	rowBalance1 := db.QueryRow("SELECT balance from cloud_pockets WHERE id = $1;", id1)
	rowBalance2 := db.QueryRow("SELECT balance from cloud_pockets WHERE id = $1;", id2)

	var balance1, balance2 float64

	rowBalance1.Scan(&balance1)
	rowBalance2.Scan(&balance2)

	byteBody, err2 := io.ReadAll(rec.Body)
	assert.NoError(t, err2)

	var t2 cloudpocket.Transfer
	err = json.Unmarshal(byteBody, &t2)

	expectedResponseBody := fmt.Sprintf(`{"id":%d,"pocket_id_source":%d,"pocket_id_target":%d,"amount":0.1}`, t2.ID, id1, id2)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, expectedResponseBody, string(byteBody))

	expectedBalance1 := 0.0
	expectedBalance2 := 0.3
	assert.Equal(t, expectedBalance1, balance1)
	assert.Equal(t, expectedBalance2, balance2)
}

func TestCreateTransferInsufficientAmountShouldBeFail(t *testing.T) {
	e := echo.New()

	cfg := config.New().All()
	db, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Error(err)
	}
	row1 := db.QueryRow("INSERT INTO cloud_pockets (name, currency, balance) VALUES ('Travel Fund', 'THB', 0.1) RETURNING id;")
	row2 := db.QueryRow("INSERT INTO cloud_pockets (name, currency, balance) VALUES ('Savings', 'THB', 0.2) RETURNING id;")

	var id1, id2 int
	row1.Scan(&id1)
	row2.Scan(&id2)

	hPocket := cloudpocket.New(db)

	e.POST("/transfers", hPocket.Transfer)

	reqBody := fmt.Sprintf(`{"pocket_id_source":%d,"pocket_id_target":%d,"amount":10}`, id1, id2)
	req := httptest.NewRequest(http.MethodPost, "/transfers", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	rowBalance1 := db.QueryRow("SELECT balance from cloud_pockets WHERE id = $1;", id1)
	rowBalance2 := db.QueryRow("SELECT balance from cloud_pockets WHERE id = $1;", id2)

	var balance1, balance2 float64

	rowBalance1.Scan(&balance1)
	rowBalance2.Scan(&balance2)

	byteBody, err2 := io.ReadAll(rec.Body)
	assert.NoError(t, err2)

	var t2 cloudpocket.Transfer
	err = json.Unmarshal(byteBody, &t2)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"message":"insufficient balance"}`, string(byteBody))
}

func TestTransferInvalidDecimalIT(t *testing.T) {
	e := echo.New()

	cfg := config.New().All()
	db, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Error(err)
	}
	row1 := db.QueryRow("INSERT INTO cloud_pockets (name, currency, balance) VALUES ('Travel Fund', 'THB', 50) RETURNING id;")
	row2 := db.QueryRow("INSERT INTO cloud_pockets (name, currency, balance) VALUES ('Savings', 'THB', 50) RETURNING id;")

	var sourceID, targetID int
	row1.Scan(&sourceID)
	row2.Scan(&targetID)

	hPocket := cloudpocket.New(db)

	e.POST("/transfers", hPocket.Transfer)

	reqBody := fmt.Sprintf(`{"pocket_id_source":%d,"pocket_id_target":%d,"amount":0.001}`, sourceID, targetID)
	req := httptest.NewRequest(http.MethodPost, "/transfers", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	byteBody, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)

	var t2 cloudpocket.Transfer
	err = json.Unmarshal(byteBody, &t2)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"message":"invalid amount"}`, string(byteBody))

	assert.Equal(t, 50.0, getBalanceFor(db, sourceID))
	assert.Equal(t, 50.0, getBalanceFor(db, targetID))
}

func getBalanceFor(db *sql.DB, id1 int) float64 {
	var rs float64
	row := db.QueryRow("SELECT balance from cloud_pockets WHERE id = $1;", id1)
	row.Scan(&rs)
	return rs
}
