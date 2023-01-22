package cloudpocket

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h handler) GetPocketById(c echo.Context) error {
	logger := mlog.L(c)

	Id := c.Param("id")
	pocket, err := GetPocketById(h.db, Id)
	if err != nil {
		logger.Error("Unsuccessful query", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Unsuccessful query", err.Error())
	}

	return c.JSON(http.StatusOK, pocket)
}

func GetPocketById(db *sql.DB, id string) (Pocket, error) {
	result := Pocket{}
	stmt, err := db.Prepare("SELECT id, name, currency, balance FROM cloud_pockets where id=$1")
	if err != nil {
		return result, errors.New("can't insert Pocket into database")
	}
	row := stmt.QueryRow(id)

	err = row.Scan(&result.ID, &result.Name, &result.Currency, &result.Balance)
	if err != nil {
		return result, errors.New("can't scan Pocket")
	}

	return result, nil
}
