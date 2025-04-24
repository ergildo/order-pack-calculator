package repositories

import (
	"context"
	"order-pack-calculator/internal/domain/entities"
)

type PackSizeRepository interface {
	Create(ctx context.Context, pack entities.PackSize) (*entities.PackSize, error)
	Update(ctx context.Context, pack entities.PackSize) error
	GetByID(ctx context.Context, ID int64) (*entities.PackSize, error)
	GetSizesByProductID(ctx context.Context, productID int64) ([]int, error)
}
