package server

import (
	"errors"
	"net/http"
	"order-pack-calculator/internal/domain/dto"
	errs "order-pack-calculator/internal/domain/errors"

	"github.com/gin-gonic/gin"
)

func ErrResponse(ctx *gin.Context, message string, err error) {
	response := dto.ErrorResponse{
		Message: message,
		Details: err.Error(),
	}
	switch {
	case errors.Is(err, errs.ErrNotFound):
		{
			ctx.JSON(http.StatusBadRequest, response)
			break
		}
	default:
		{
			ctx.JSON(http.StatusInternalServerError, response)
			break
		}
	}
}
