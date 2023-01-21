package pocket

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	cStmt = "INSERT INTO transfers (pocket_id_source, pocket_id_target, amount) VALUES ($1, $2, $3) RETURNING id;"
)

func (h handler) Transfer(c echo.Context) error {
	logger := mlog.L(c)
	var transfer TransferModel
	ctx := c.Request().Context()
	err := c.Bind(&transfer)
	if transfer.Amount < 1 {
		logger.Error("bad request invalid amount")
		return echo.NewHTTPError(http.StatusBadRequest, "bad request invalid amount")
	}
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
	}
	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt, transfer.PocketIDSource, transfer.PocketIDTarget, transfer.Amount).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}
	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	transfer.ID = lastInsertId
	return c.JSON(http.StatusCreated, transfer)
}
