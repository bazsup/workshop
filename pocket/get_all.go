package pocket

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetAll(c echo.Context) error {
	stml, err := h.db.Prepare("SELECT * FROM cloud_pockets")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	rows, err := stml.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Can't query")
	}
	pockets := []PocketModel{}
	for rows.Next() {
		p := PocketModel{}
		err := rows.Scan(&p.ID, &p.Name, &p.Currency, &p.Balance)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Can't scan pocket cloud")
		}

		pockets = append(pockets, p)
	}

	return c.JSON(http.StatusOK, pockets)
}
