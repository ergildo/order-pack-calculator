package server

import (
	"net/http"
	"order-pack-calculator/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// CalculatePackSizeHandler godoc
// @Summary      Calculate optimal pack sizes
// @Description  Calculates the optimal pack sizes for a given order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order  body      dto.CalculatePackSizesRequest  true  "Order details"
// @Success      200    {object}  dto.OptimalPackSizesResponse
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Router       /api/v1/orders/calculate [post]
func (s *Server) CalculatePackSizeHandler(ctx *gin.Context) {
	var order dto.CalculatePackSizesRequest
	err := ctx.BindJSON(&order)
	if err != nil {
		ErrResponse(ctx, "unable to parse request", err)
		return
	}

	optimal, err := s.packSizeService.CalcOptimalPacks(ctx, order)

	if err != nil {
		ErrResponse(ctx, "unable to calculate pack sizes", err)
		return
	}
	ctx.JSON(http.StatusOK, optimal)
}
