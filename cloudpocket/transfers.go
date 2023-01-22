package cloudpocket

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	cStmt = "INSERT INTO transfers (pocket_id_source, pocket_id_target, amount) VALUES ($1, $2, $3) RETURNING id;"
	uStmt = "UPDATE cloud_pockets SET balance=$1 WHERE id=$2;"
)

func (h handler) Transfer(c echo.Context) error {
	logger := mlog.L(c)
	var t Transfer
	ctx := c.Request().Context()
	err := c.Bind(&t)
	if !t.IsValidAmount() {
		logger.Error("bad request invalid amount")
		return c.JSON(http.StatusBadRequest, echo.NewHTTPError(http.StatusBadRequest, "invalid amount"))
	}
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.NewHTTPError(http.StatusBadRequest, "bad request body"))
	}

	var sPocket, tPocket Pocket
	sPocket, err = GetPocketById(h.db, t.PocketIDSource)
	if err != nil {
		logger.Error("get pocket error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
	}
	tPocket, err = GetPocketById(h.db, t.PocketIDTarget)
	if err != nil {
		logger.Error("get pocket error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
	}

	_, err = h.db.Exec(uStmt, round(sPocket.Balance-t.Amount), sPocket.ID)
	if err != nil {
		logger.Error("update source balance error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
	}

	if sPocket.Balance < t.Amount {
		logger.Error("insufficient balance")
		return c.JSON(http.StatusBadRequest, echo.NewHTTPError(http.StatusBadRequest, "insufficient balance"))
	}

	_, err = h.db.Exec(uStmt, round(tPocket.Balance+t.Amount), tPocket.ID)
	if err != nil {
		logger.Error("update target balance error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
	}

	var lastInsertId int
	err = h.db.QueryRowContext(ctx, cStmt, t.PocketIDSource, t.PocketIDTarget, t.Amount).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}
	logger.Info("create successfully", zap.Int("id", lastInsertId))
	t.ID = lastInsertId
	return c.JSON(http.StatusCreated, t)
}
