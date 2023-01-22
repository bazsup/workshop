//go:build unit

package cloudpocket_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/cloudpocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransferShouldBeSuccess(t *testing.T) {
	t.Run("Create transfer should be succesfully", func(t *testing.T) {

		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
		assert.NoError(t, err)
		defer db.Close()

		mQueryRows1 := sqlmock.NewRows([]string{"id", "name", "currency", "balance"}).
			AddRow(1, "Travel", "THB", 0.01)

		mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, name, currency, balance FROM cloud_pockets where id=$1")).
			ExpectQuery().
			WithArgs(1).
			WillReturnRows(mQueryRows1)

		mQueryRows2 := sqlmock.NewRows([]string{"id", "name", "currency", "balance"}).
			AddRow(2, "Shoping", "THB", 0.02)
		mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, name, currency, balance FROM cloud_pockets where id=$1")).
			ExpectQuery().
			WithArgs(2).
			WillReturnRows(mQueryRows2)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE cloud_pockets SET balance=$1 WHERE id=$2;")).
			WithArgs(0.0, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta("UPDATE cloud_pockets SET balance=$1 WHERE id=$2;")).
			WithArgs(0.03, 2).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mIntsertRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO transfers (pocket_id_source, pocket_id_target, amount) VALUES ($1, $2, $3) RETURNING id;")).
			WithArgs(1, 2, 0.01).
			WillReturnRows(mIntsertRows)

		reqBody := `{"pocket_id_source":1,"pocket_id_target":2,"amount":0.01}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/cloud-pocket/transfers", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		h := cloudpocket.New(db)

		err = h.Transfer(ctx)

		if err != nil {
			t.Errorf("should not return error but it got %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("should status bed request but it got %v", ctx.Response().Status)
			byteBuffer, _ := io.ReadAll(rec.Body)
			t.Error(string(byteBuffer))
		}

	})
}

func TestCreateTransferInvalidAmountShouldBeFail(t *testing.T) {
	t.Skip()
	t.Run("Create transfer with invalid amount should got error", func(t *testing.T) {

		reqBody := `{"pocket_id_source": 1,"pocket_id_target": 1,"amount": 0.001}`

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/cloud-pocket/transfers", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		h := cloudpocket.New(db)

		err = h.Transfer(ctx)

		if err != nil {
			t.Errorf("should not return error but it got %v", err)
		}
		if ctx.Response().Status != http.StatusBadRequest {
			t.Errorf("should status bed request but it got %v", ctx.Response().Status)
		}
		assert.JSONEq(t, `{"message":"invalid amount"}`, rec.Body.String())
	})
}
