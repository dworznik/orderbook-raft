package fsm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dworznik/orderbook"
	"github.com/hashicorp/raft"
	"github.com/shopspring/decimal"
)

type Command struct {
	Op    string
	Order orderbook.Order
}

func (c *Command) UnmarshalJSON(data []byte) error {
	cmd := struct {
		Op    string          `json:"op"`
		Order orderbook.Order `json:"order"`
	}{}

	if err := json.Unmarshal(data, &cmd); err != nil {
		return err
	}

	if cmd.Op != "MARKET" && cmd.Op != "LIMIT" && cmd.Op != "CANCEL" {
		return errors.New("invalid op")
	}

	c.Op = cmd.Op
	c.Order = cmd.Order

	return nil
}

func (c *Command) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Op    string          `json:"op"`
			Order orderbook.Order `json:"order"`
		}{
			Op:    c.Op,
			Order: c.Order,
		},
	)
}

type MarketResult struct {
	done                     []*orderbook.Order
	partial                  *orderbook.Order
	partialQuantityProcessed decimal.Decimal
	quantityLeft             decimal.Decimal
}

type LimitResult struct {
	Done                     []*orderbook.Order `json:"done"`
	Partial                  *orderbook.Order   `json:"partial"`
	PartialQuantityProcessed decimal.Decimal    `json:"partialQuantityProcessed"`
}

type CancelResult struct {
	Order *orderbook.Order `json:"order"`
}

type ApplyResponse struct {
	Error error
	Data  interface{}
}

type orderbookFSM struct {
	orderbook.OrderBook
}

func (ob *orderbookFSM) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		cmd := Command{}
		if err := json.Unmarshal(log.Data, &cmd); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error marshalling command %s\n", err.Error())
			return nil
		}

		op := strings.ToUpper(strings.TrimSpace(cmd.Op))

		switch op {
		case "MARKET":
			o := cmd.Order
			done, partial, partialQuantityProcessed, quantityLeft, err := ob.ProcessMarketOrder(o.Side(), o.Quantity())
			return &ApplyResponse{
				Error: err,
				Data:  MarketResult{done, partial, partialQuantityProcessed, quantityLeft},
			}
		case "LIMIT":
			o := cmd.Order
			done, partial, partialQuantityProcessed, err := ob.ProcessLimitOrder(o.Side(), o.ID(), o.Quantity(), o.Price())
			return &ApplyResponse{
				Error: err,
				Data:  LimitResult{done, partial, partialQuantityProcessed},
			}
		case "CANCEL":
			o := cmd.Order
			order := ob.CancelOrder(o.ID())
			if order != nil {
				return &ApplyResponse{
					Error: nil,
					Data:  CancelResult{order},
				}
			} else {
				return &ApplyResponse{
					Error: fmt.Errorf("order not found"),
					Data:  nil,
				}
			}

		}
	}

	_, _ = fmt.Fprintf(os.Stderr, "not raft log command type\n")
	return nil
}

func (ob *orderbookFSM) Snapshot() (raft.FSMSnapshot, error) {
	return newSnapshotNoop()
}

func (ob *orderbookFSM) Restore(rClose io.ReadCloser) error {
	return nil
}

func NewOrderbookFSM(ob orderbook.OrderBook) raft.FSM {
	return &orderbookFSM{
		ob,
	}
}
