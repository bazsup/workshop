//go:build unit

package cloudpocket_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/cloudpocket"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransferShouldBeSuccess(t *testing.T) {
	t.Skip()
	tests := []struct {
		name       string
		cfgFlag    config.FeatureFlag
		sqlFn      func() (*sql.DB, error)
		reqBody    string
		wantStatus int
		wantBody   string
	}{
		{"create pocket succesfully",
			config.FeatureFlag{},
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				row := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`INSERT INTO transfers (pocket_id_source, pocket_id_target, amount) VALUES ($1, $2, $3) RETURNING id;`).WithArgs(1, 3, 1.0).WillReturnRows(row)
				return db, err
			},
			`{"pocket_id_source": 1,"pocket_id_target": 3,"amount": 1.0}`,
			http.StatusCreated,
			`{"id":1,"pocket_id_source":1,"pocket_id_target":3,"amount":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/cloud-pocket/transfers", strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			db, err := tc.sqlFn()
			h := cloudpocket.New(db)
			// Assertions
			assert.NoError(t, err)
			if assert.NoError(t, h.Transfer(c)) {
				assert.Equal(t, tc.wantStatus, rec.Code)
				assert.JSONEq(t, tc.wantBody, rec.Body.String())
			}
		})
	}
}

func TestCreateTransferInvalidAmountShouldBeFail(t *testing.T) {
	t.Run("Create transfer with invalid amount should got error", func(t *testing.T) {

		reqBody := `{"pocket_id_source": 1,"pocket_id_target": 1,"amount": 0.001}`
		wantStatus := http.StatusBadRequest
		wantBody := `{"message":"invalid amount"}`

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
		if ctx.Response().Status != wantStatus {
			t.Errorf("should status bed request but it got %v", ctx.Response().Status)
		}
		assert.JSONEq(t, wantBody, rec.Body.String())
	})
}
