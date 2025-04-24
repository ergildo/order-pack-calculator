package server

import (
	"net/http"
	"order-pack-calculator/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// UpdatePackSizeHandler godoc
// @Summary      Update pack sizes
// @Description  Updates existing pack sizes
// @Tags         packsizes
// @Accept       json
// @Produce      json
// @Param        packSize  body      dto.UpdatePackSizeRequest  true  "Updated pack size details"
// @Success      200       "OK"
// @Failure      400       {object}  dto.ErrorResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /api/v1/packsizes [patch]
func (s *Server) UpdatePackSizeHandler(ctx *gin.Context) {
	var request dto.UpdatePackSizeRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ErrResponse(ctx, "unable to parse request", err)
		return
	}

	err = s.packSizeService.Update(ctx, request)

	if err != nil {
		ErrResponse(ctx, "unable to update pack sizes", err)
		return
	}
	ctx.Status(http.StatusOK)
}
