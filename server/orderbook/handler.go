package orderbook_handler

import (
	"github.com/dworznik/orderbook"
	"github.com/hashicorp/raft"
)

type handler struct {
	raft *raft.Raft
	ob   *orderbook.OrderBook
}

type responseBase struct {
	Status      string `json:"status"`
	ErrorCode   string `json:"error_code,omitempty"`
	TimeSpentMs int64  `json:"time_spent_ms"`
}

func New(raft *raft.Raft, ob *orderbook.OrderBook) *handler {
	return &handler{
		raft: raft,
		ob:   ob,
	}
}
