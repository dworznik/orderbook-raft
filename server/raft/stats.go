package raft_handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// StatsRaftHandler get raft status
func (h handler) StatsRaftHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Here is the raft status",
		"data":    h.raft.Stats(),
	})
}
