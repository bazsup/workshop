//go:build unit

package pocket

import(
	"github.com/DATA-DOG/go-sqlmock"
)

func TestgetPocketById(T *testing.T) {
	pocket := PocketModel{
		ID = 1,
		Name = "shoping"
		Currency = "THB",
		Balance = 100.0

	}
	db,mock,_ := sqlmock.New()
    row := sqlmock.NewRows([]string{"id, name, currency, balance"}).AddRow(pocket)


}
