package orderbook_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dworznik/orderbook"
	"github.com/dworznik/orderbook-raft/fsm"
	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type requestLimit struct {
	Side     orderbook.Side  `json:"side"`
	Price    decimal.Decimal `json:"price" validate:"dgt=0"`
	Quantity decimal.Decimal `json:"quantity" validate:"dgt=0"`
}

type dataLimit struct {
	Order  orderbook.Order `json:"order"`
	Result fsm.LimitResult `json:"result"`
}

type responseLimit struct {
	responseBase
	Data dataLimit `json:"data"`
}

type dataCancel struct {
	Result fsm.CancelResult `json:"result"`
}

type responseCancel struct {
	responseBase
	Data dataCancel `json:"data"`
}

type requestCancel struct {
	OrderId string `json:"orderId"`
}

func (h handler) Limit(c echo.Context) error {
	startTime := time.Now()

	var form = requestLimit{}

	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, responseBase{"error", fmt.Sprintf("invalid input: %s", err.Error()), time.Since(startTime).Milliseconds()})
	}

	if err := c.Validate(&form); err != nil {
		return c.JSON(http.StatusBadRequest, responseBase{"error", fmt.Sprintf("validation error: %s", err.Error()), time.Since(startTime).Milliseconds()})
	}

	if h.raft.State() != raft.Leader {
		return c.JSON(http.StatusUnprocessableEntity, responseBase{"error", "not leader", time.Since(startTime).Milliseconds()})
	}

	order := orderbook.NewOrder(uuid.New().String(), form.Side, form.Quantity, form.Price, time.Now())
	payload := fsm.Command{
		Op:    "LIMIT",
		Order: *order,
	}

	data, err := json.Marshal(&payload)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error preparing saving data payload: %s", err.Error()),
		})
	}

	applyFuture := h.raft.Apply(data, 500*time.Millisecond)
	if err := applyFuture.Error(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error persisting data in raft cluster: %s", err.Error()),
		})
	}

	res, ok := applyFuture.Response().(*fsm.ApplyResponse)

	if !ok {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": "error response is not match apply response",
		})
	}

	if res.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": res.Error.Error(),
		})
	}

	result, _ := res.Data.(fsm.LimitResult)
	resp := &responseLimit{
		responseBase: responseBase{"success", "", time.Since(startTime).Milliseconds()},
		Data:         dataLimit{*order, result},
	}
	return c.JSON(http.StatusOK, resp)
}

func (h handler) Cancel(c echo.Context) error {
	startTime := time.Now()

	var form = requestCancel{}

	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, responseBase{"error", fmt.Sprintf("invalid input: %s", err.Error()), time.Since(startTime).Milliseconds()})
	}

	if err := c.Validate(&form); err != nil {
		return c.JSON(http.StatusBadRequest, responseBase{"error", fmt.Sprintf("validation error: %s", err.Error()), time.Since(startTime).Milliseconds()})
	}

	if h.raft.State() != raft.Leader {
		return c.JSON(http.StatusUnprocessableEntity, responseBase{"error", "not leader", time.Since(startTime).Milliseconds()})
	}

	order := orderbook.NewOrder(form.OrderId, 0, decimal.Zero, decimal.Zero, time.Now())
	payload := fsm.Command{
		Op:    "CANCEL",
		Order: *order,
	}

	data, err := json.Marshal(&payload)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error preparing saving data payload: %s", err.Error()),
		})
	}

	applyFuture := h.raft.Apply(data, 500*time.Millisecond)
	if err := applyFuture.Error(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error persisting data in raft cluster: %s", err.Error()),
		})
	}

	res, ok := applyFuture.Response().(*fsm.ApplyResponse)

	if !ok {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": "error response is not match apply response",
		})
	}

	if res.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": res.Error.Error(),
		})
	}

	result, _ := res.Data.(fsm.CancelResult)
	resp := &responseCancel{
		responseBase: responseBase{"success", "", time.Since(startTime).Milliseconds()},
		Data:         dataCancel{result},
	}
	return c.JSON(http.StatusOK, resp)
}
