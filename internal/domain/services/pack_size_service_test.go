package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"order-pack-calculator/internal/domain/dto"
	"order-pack-calculator/internal/domain/entities"

	"order-pack-calculator/mocks"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockPackSizeRepository(ctrl)
	service := NewPackSizeService(repo)

	t.Run("success", func(t *testing.T) {
		req := dto.CreatePackSizeRequest{ProductID: 1, Size: 10}
		saved := entities.PackSize{ID: 100, ProductID: 1, Size: 10, Active: true}

		ctx := context.Background()

		repo.EXPECT().Create(ctx, entities.PackSize{ProductID: 1, Size: 10}).Return(&saved, nil)

		resp, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, int64(100), resp.ID)
		assert.Equal(t, 10, resp.Size)
	})

	t.Run("repository error", func(t *testing.T) {
		repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("repo error"))
		_, err := service.Create(context.Background(), dto.CreatePackSizeRequest{})
		assert.Error(t, err)
	})
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockPackSizeRepository(ctrl)
	service := NewPackSizeService(repo)

	t.Run("update size and active", func(t *testing.T) {
		newSize := 20
		newActive := false
		existing := &entities.PackSize{ID: 1, ProductID: 1, Size: 10, Active: true}
		updated := *existing
		updated.Size = newSize
		updated.Active = newActive

		repo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(existing, nil)
		repo.EXPECT().Update(gomock.Any(), updated).Return(nil)

		err := service.Update(context.Background(), dto.UpdatePackSizeRequest{
			ID:     1,
			Size:   &newSize,
			Active: &newActive,
		})
		assert.NoError(t, err)
	})

	t.Run("get by id error", func(t *testing.T) {
		repo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(nil, errors.New("not found"))
		err := service.Update(context.Background(), dto.UpdatePackSizeRequest{ID: 1})
		assert.Error(t, err)
	})

	t.Run("update error", func(t *testing.T) {
		pack := &entities.PackSize{ID: 1, ProductID: 1, Size: 10, Active: true}
		repo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(pack, nil)
		repo.EXPECT().Update(gomock.Any(), *pack).Return(errors.New("update failed"))
		err := service.Update(context.Background(), dto.UpdatePackSizeRequest{ID: 1})
		assert.Error(t, err)
	})
}

func TestCalcOptimalPacks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockPackSizeRepository(ctrl)
	service := NewPackSizeService(repo)

	tests := []struct {
		name       string
		packs      []int
		orderQty   int
		expectResp *dto.OptimalPackSizesResponse
	}{
		{
			"desired output",
			[]int{23, 31, 53},
			500000,
			&dto.OptimalPackSizesResponse{
				PackCombination: []dto.PackDetail{{Size: 23, Count: 2}, {Size: 31, Count: 7}, {Size: 53, Count: 9429}},
				TotalItems:      500000,
				TotalPacks:      9438,
			},
		},

		{
			"exact match",
			[]int{5, 10, 20},
			20,
			&dto.OptimalPackSizesResponse{
				PackCombination: []dto.PackDetail{{Size: 20, Count: 1}},
				TotalItems:      20,
				TotalPacks:      1,
			},
		},
		{
			"multiple combinations",
			[]int{3, 7},
			10,
			&dto.OptimalPackSizesResponse{
				PackCombination: []dto.PackDetail{{Size: 3, Count: 1}, {Size: 7, Count: 1}},
				TotalItems:      10,
				TotalPacks:      2,
			},
		},
		{
			"overfill minimal",
			[]int{6, 8},
			10,
			&dto.OptimalPackSizesResponse{
				PackCombination: []dto.PackDetail{{Size: 6, Count: 2}},
				TotalItems:      12,
				TotalPacks:      2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().GetSizesByProductID(gomock.Any(), gomock.Any()).Return(tt.packs, nil)
			resp, err := service.CalcOptimalPacks(context.Background(), dto.CalculatePackSizesRequest{ProductID: 1, OrderQuantity: tt.orderQty})
			assert.NoError(t, err)
			assert.Equal(t, tt.expectResp.TotalItems, resp.TotalItems)
			assert.Equal(t, tt.expectResp.TotalPacks, resp.TotalPacks)
			assert.ElementsMatch(t, tt.expectResp.PackCombination, resp.PackCombination)
		})
	}

	t.Run("repository error", func(t *testing.T) {
		repo.EXPECT().GetSizesByProductID(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
		_, err := service.CalcOptimalPacks(context.Background(), dto.CalculatePackSizesRequest{ProductID: 1, OrderQuantity: 10})
		assert.Error(t, err)
	})
}
