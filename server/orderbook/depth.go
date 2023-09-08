package orderbook_handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dworznik/orderbook"

	"github.com/labstack/echo/v4"
)

type dataDepth struct {
	Asks []*orderbook.PriceLevel `json:"asks"`
	Bids []*orderbook.PriceLevel `json:"bids"`
}

type responseDepth struct {
	responseBase
	Data dataDepth `json:"data"`
}

func (h handler) Depth(c echo.Context) error {
	startTime := time.Now()

	form := requestLimit{}
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error binding: %s", err.Error()),
		})
	}

	asks, bids := h.ob.Depth()
	resp := &responseDepth{
		responseBase: responseBase{"success", "", time.Since(startTime).Milliseconds()},
		Data:         dataDepth{asks, bids},
	}
	return c.JSON(http.StatusOK, resp)
}
