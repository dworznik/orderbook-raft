package server

import (
	"net/http"
	_ "net/http/pprof"
	"reflect"
	"time"

	"github.com/dworznik/orderbook"
	orderbook_handler "github.com/dworznik/orderbook-raft/server/orderbook"
	raft_handler "github.com/dworznik/orderbook-raft/server/raft"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/raft"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shopspring/decimal"
)

// srv struct handling server
type srv struct {
	listenAddress string
	raft          *raft.Raft
	echo          *echo.Echo
}

// Start start the server
func (s srv) Start() error {
	return s.echo.StartServer(&http.Server{
		Addr:         s.listenAddress,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// https://github.com/go-playground/validator/issues/515#issuecomment-1235164600
func decimalGreaterThan(validate *validator.Validate) error {
	validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if valuer, ok := field.Interface().(decimal.Decimal); ok {
			return valuer.String()
		}
		return nil
	}, decimal.Decimal{})
	if err := validate.RegisterValidation("dgt", func(fl validator.FieldLevel) bool {
		data, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		value, err := decimal.NewFromString(data)
		if err != nil {
			return false
		}
		baseValue, err := decimal.NewFromString(fl.Param())
		if err != nil {
			return false
		}
		return value.GreaterThan(baseValue)
	}); err != nil {
		return err
	}

	return nil
}

func New(listenAddr string, ob *orderbook.OrderBook, r *raft.Raft) *srv {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Pre(middleware.RemoveTrailingSlash())

	validator := validator.New()
	decimalGreaterThan(validator)
	e.Validator = &CustomValidator{validator: validator}

	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	raftHandler := raft_handler.New(r)
	e.POST("/raft/join", raftHandler.JoinRaftHandler)
	e.GET("/raft/stats", raftHandler.StatsRaftHandler)

	orderBookHandler := orderbook_handler.New(r, ob)

	e.POST("/limit", orderBookHandler.Limit)
	e.POST("/cancel", orderBookHandler.Cancel)
	e.GET("/depth", orderBookHandler.Depth)

	return &srv{
		listenAddress: listenAddr,
		echo:          e,
		raft:          r,
	}
}
