package services

import (
	"context"
	"order-pack-calculator/internal/domain/dto"
)

type PackSizeService interface {
	CalcOptimalPacks(context.Context, dto.CalculatePackSizesRequest) (*dto.OptimalPackSizesResponse, error)
	Create(context.Context, dto.CreatePackSizeRequest) (*dto.PackSizeResponse, error)
	Update(context.Context, dto.UpdatePackSizeRequest) error
	GetAll(ctx context.Context) ([]dto.PackSizeResponse, error)
}
