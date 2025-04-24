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

func TestUpdatePackSizeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockPackSizeService(ctrl)
		s := &Server{packSizeService: mockService}

		size := 15
		active := true
		reqBody := dto.UpdatePackSizeRequest{ID: 1, Size: &size, Active: &active}
		mockService.EXPECT().Update(gomock.Any(), reqBody).Return(nil)

		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/packsizes", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r, _ := gin.CreateTestContext(w)
		r.Request = req

		s.UpdatePackSizeHandler(r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		s := &Server{}

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/packsizes", bytes.NewBuffer([]byte(`invalid`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r, _ := gin.CreateTestContext(w)
		r.Request = req

		s.UpdatePackSizeHandler(r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error - service failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockPackSizeService(ctrl)
		s := &Server{packSizeService: mockService}

		size := 15
		reqBody := dto.UpdatePackSizeRequest{ID: 1, Size: &size}
		mockService.EXPECT().Update(gomock.Any(), reqBody).Return(errors.New("update failed"))

		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/packsizes", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r, _ := gin.CreateTestContext(w)
		r.Request = req

		s.UpdatePackSizeHandler(r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
