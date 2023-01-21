package pocket

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (handler) GetAll(c echo.Context) error {
	pockets := []PocketModel{
		{
			ID: 12345,
			Name: "Travel Fund",
			Currency: "THB",
			Balance: 100,
		},
		{
			ID: 67890,
			Name: "Savings",
			Currency: "THB",
			Balance: 200,
		},
	}
	return c.JSON(http.StatusOK, pockets)
}
