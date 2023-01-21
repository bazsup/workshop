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
			Category: "Vacation",
			Currency: "THB",
			Balance: 100,
		},
		{
			ID: 67890,
			Name: "Savings",
			Category: "Emergency Fund",
			Currency: "THB",
			Balance: 200,
		},
	}
	return c.JSON(http.StatusOK, pockets)
}
