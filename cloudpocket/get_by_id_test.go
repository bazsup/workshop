//go:build unit

package cloudpocket_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/cloudpocket"
	"github.com/stretchr/testify/assert"
)

func TestGetPocketById(t *testing.T) {

	db, mock, _ := sqlmock.New()
	row := sqlmock.NewRows([]string{"id", "name", "currency", "balance"}).AddRow("1", "shoping", "THB", "100.0")
	mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, name, currency, balance FROM cloud_pockets where id=$1")).ExpectQuery().WithArgs(1).WillReturnRows(row)

	result, err := cloudpocket.GetPocketById(db, 1)

	assert.Nil(t, err)
	assert.EqualValues(t, result.ID, 1)
	assert.EqualValues(t, result.Name, "shoping")
	assert.EqualValues(t, result.Currency, "THB")
	assert.EqualValues(t, result.Balance, 100.0)

}
