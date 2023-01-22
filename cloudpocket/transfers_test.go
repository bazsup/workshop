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

func TestCreateTransfer(t *testing.T) {
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
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.reqBody))
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
