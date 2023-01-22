//go:build unit

package cloudpocket_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/cloudpocket"
	"github.com/stretchr/testify/assert"
)

func TestgetPocketById(t *testing.T) {
	pocket := cloudpocket.Pocket{
		ID:       1,
		Name:     "shoping",
		Currency: "THB",
		Balance:  100.0,
	}
	db, mock, _ := sqlmock.New()
	row := sqlmock.NewRows([]string{"id, name, currency, balance"}).AddRow(pocket)
	mock.ExpectPrepare("SELECT ID, Name, Currency, Balance FROM cloud_pockets").ExpectQuery().WillReturnRows(row)

	result, err := cloudpocket.GetPocketById(db, "1")
	assert.Nil(t, err)
	assert.EqualValues(t, result.ID, pocket.ID)
	assert.EqualValues(t, result.Name, pocket.Name)
	assert.EqualValues(t, result.Currency, pocket.Currency)
	assert.EqualValues(t, result.Balance, pocket.Balance)

}
