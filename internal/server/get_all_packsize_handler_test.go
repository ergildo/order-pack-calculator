package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"order-pack-calculator/internal/domain/dto"
	"order-pack-calculator/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPackSizeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockPackSizeService(ctrl)
		s := &Server{packSizeService: mockService}

	
		respBody := []dto.PackSizeResponse{{ID: 1, ProductID: 1, Size: 10, Active: true}}

		mockService.EXPECT().GetAll(gomock.Any()).Return(respBody, nil)

	
		req := httptest.NewRequest(http.MethodGet, "/api/v1/packsizes", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r, _ := gin.CreateTestContext(w)
		r.Request = req

		s.GetAllPackSizeHandler(r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("internal server error - service failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockPackSizeService(ctrl)
		s := &Server{packSizeService: mockService}

		
		mockService.EXPECT().GetAll(gomock.Any()).Return(nil, errors.New("db error"))

	
		req := httptest.NewRequest(http.MethodGet, "/api/v1/packsizes", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r, _ := gin.CreateTestContext(w)
		r.Request = req

		s.GetAllPackSizeHandler(r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
