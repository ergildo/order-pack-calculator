package server

import (
	"net/http"
	"order-pack-calculator/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// CreatePackSizeHandler godoc
// @Summary      Create pack sizes
// @Description  Creates new pack sizes
// @Tags         packsizes
// @Accept       json
// @Produce      json
// @Param        packSize  body      dto.CreatePackSizeRequest  true  "Pack size details"
// @Success      200       {object}  dto.PackSizeResponse
// @Failure      400       {object}  dto.ErrorResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /api/v1/packsizes [post]
func (s *Server) CreatePackSizeHandler(ctx *gin.Context) {
	var request dto.CreatePackSizeRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ErrResponse(ctx, "unable to parse request", err)
		return
	}

	response, err := s.packSizeService.Create(ctx, request)

	if err != nil {
		ErrResponse(ctx, "unable to create pack sizes", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
