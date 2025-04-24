package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePackSizeHandler godoc
// @Summary      Get All pack sizes
// @Description  Get All pack sizes
// @Tags         packsizes
// @Accept       json
// @Produce      json
// @Success      200       {array}  dto.PackSizeResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /api/v1/packsizes [get]
func (s *Server) GetAllPackSizeHandler(ctx *gin.Context) {
	

	response, err := s.packSizeService.GetAll(ctx)

	if err != nil {
		ErrResponse(ctx, "unable to get pack sizes", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
