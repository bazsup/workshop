package cloudpocket

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h handler) CreateSubPocket(c echo.Context) error {
	logger := mlog.L(c)
	pocket := new(Pocket)
	ctx := c.Request().Context()

	err := c.Bind(&pocket)

	if err != nil {
		return c.JSON(http.StatusBadRequest, pocket)
	}

	var lastInsertId int

	err = h.db.QueryRowContext(ctx, `INSERT INTO cloud_pockets (name,currency,balance) VALUES ($1,$2,$3) RETURNING id`, pocket.Name, pocket.Currency, pocket.Balance).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}

	pocket.ID = lastInsertId

	return c.JSON(http.StatusCreated, pocket)

}
