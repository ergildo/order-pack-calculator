package server

import (
	"bytes"
	"encoding/json"
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

func TestCalculatePackSizeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockPackSizeService(ctrl)
		s := &Server{packSizeService: mockService}

		reqBody := dto.CalculatePackSizesRequest{ProductID: 1, OrderQuantity: 20}
		respBody := &dto.OptimalPackSizesResponse{
			PackCombination: []dto.PackDetail{{Size: 10, Count: 2}},
			TotalItems:      20,
			TotalPacks:      2,
		}

		mockService.EXPECT().CalcOptimalPacks(gomock.Any(), reqBody).Return(respBody, nil)

		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/calculate", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r, _ := gin.CreateTestContext(w)
		r.Request = req

		s.CalculatePackSizeHandler(r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		s := &Server{}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/calculate", bytes.NewBuffer([]byte(`invalid`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r, _ := gin.CreateTestContext(w)
		r.Request = req

		s.CalculatePackSizeHandler(r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error - service failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockPackSizeService(ctrl)
		s := &Server{packSizeService: mockService}

		reqBody := dto.CalculatePackSizesRequest{ProductID: 1, OrderQuantity: 10}
		mockService.EXPECT().CalcOptimalPacks(gomock.Any(), reqBody).Return(nil, errors.New("calculation failed"))

		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/calculate", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r, _ := gin.CreateTestContext(w)
		r.Request = req

		s.CalculatePackSizeHandler(r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
