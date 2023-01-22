package cloudpocket

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h handler) GetPocketById(c echo.Context) error {
	logger := mlog.L(c)

	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	pocket, err := GetPocketById(h.db, ID)
	if err != nil {
		logger.Error("Unsuccessful query", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Unsuccessful query", err.Error())
	}

	return c.JSON(http.StatusOK, pocket)
}

func GetPocketById(db *sql.DB, id int) (Pocket, error) {
	result := Pocket{}
	stmt, err := db.Prepare("SELECT id, name, currency, balance FROM cloud_pockets where id=$1")
	if err != nil {
		fmt.Println(err)
		return result, errors.New("can't Read Pocket from database")
	}
	row := stmt.QueryRow(id)

	err = row.Scan(&result.ID, &result.Name, &result.Currency, &result.Balance)
	if err != nil {
		fmt.Println(err)
		return result, errors.New("can't scan Pocket")
	}

	return result, nil
}
