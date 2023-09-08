package fsm

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dworznik/orderbook"
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
)

func Test(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Orderbook", func() {
		g.It("Should unmarshall limit order message", func() {
			cmd := &Command{}

			err := json.Unmarshal([]byte(`{ "op": "LIMIT", "order": { "id": "id", "side": "buy", "timestamp": "2019-10-10T10:10:10Z", "quantity": "10.0", "price": "120.20" } }`), &cmd)
			Ω(err).Should(BeNil())
			fmt.Printf("%+v", cmd)
			Ω(cmd.Op).Should(Equal("LIMIT"))
			Ω(cmd.Order.ID()).Should(Equal("id"))
			Ω(cmd.Order.Side()).Should(Equal(orderbook.Buy))
			Ω(cmd.Order.Quantity().Cmp(decimal.NewFromFloat(10.0))).Should(Equal(0))
			Ω(cmd.Order.Price().Cmp(decimal.NewFromFloat(120.20))).Should(Equal(0))
		})

		g.It("Should unmarshall market order message", func() {
			cmd := &Command{}

			err := json.Unmarshal([]byte(`{ "op": "MARKET",  "order": { "side": "buy", "timestamp": "2019-10-10T10:10:10Z", "quantity": "10.0", "price": "120.20" }}`), &cmd)
			Ω(err).Should(BeNil())
			Ω(cmd.Op).Should(Equal("MARKET"))
			Ω(cmd.Order.Side()).Should(Equal(orderbook.Buy))
			Ω(cmd.Order.Quantity().Cmp(decimal.NewFromFloat(10.0))).Should(Equal(0))
			Ω(cmd.Order.Price().Cmp(decimal.NewFromFloat(120.20))).Should(Equal(0))
		})

		g.It("Should unmarshall cancel order message", func() {
			cmd := &Command{}

			err := json.Unmarshal([]byte(`{ "op": "CANCEL", "order": { "id": "id"  }}`), &cmd)
			Ω(err).Should(BeNil())
			Ω(cmd.Op).Should(Equal("CANCEL"))
			Ω(cmd.Order.ID()).Should(Equal("id")) // Ignore the rest
		})

		g.It("Should unmarshall unknown command", func() {
			cmd := &Command{}
			err := json.Unmarshal([]byte(`{ "op": "UNKNOWN", "id": "id"  }`), &cmd)
			Ω(err.Error()).Should(Equal("invalid op"))
		})
	})
}
