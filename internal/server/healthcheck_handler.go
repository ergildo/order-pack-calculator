package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthHandler godoc
// @Summary      Health check
// @Description  Returns the health status of the application
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/health [get]
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.dbService.Health())
}
