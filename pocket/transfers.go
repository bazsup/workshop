package pocket

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var ts = []TransferModel{}

func (h handler) Transfer(c echo.Context) error {
	logger := mlog.L(c)
	var transfer TransferModel
	err := c.Bind(&transfer)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
	}

	ts = append(ts, transfer)
	if ts == nil {
		logger.Error("insert row error", zap.Error(err))
		return err
	}

	return c.JSON(http.StatusCreated, transfer)
}
