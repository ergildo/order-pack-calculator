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

func TestCreatePackSizeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockPackSizeService(ctrl)
		s := &Server{packSizeService: mockService}

		reqBody := dto.CreatePackSizeRequest{ProductID: 1, Size: 10}
		respBody := &dto.PackSizeResponse{ID: 1, ProductID: 1, Size: 10, Active: true}

		mockService.EXPECT().Create(gomock.Any(), reqBody).Return(respBody, nil)

		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/packsizes", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r, _ := gin.CreateTestContext(w)
		r.Request = req

		s.CreatePackSizeHandler(r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		s := &Server{}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/packsizes", bytes.NewBuffer([]byte(`invalid json`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r, _ := gin.CreateTestContext(w)
		r.Request = req

		s.CreatePackSizeHandler(r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error - service failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockPackSizeService(ctrl)
		s := &Server{packSizeService: mockService}

		reqBody := dto.CreatePackSizeRequest{ProductID: 1, Size: 10}
		mockService.EXPECT().Create(gomock.Any(), reqBody).Return(nil, errors.New("db error"))

		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/packsizes", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r, _ := gin.CreateTestContext(w)
		r.Request = req

		s.CreatePackSizeHandler(r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
