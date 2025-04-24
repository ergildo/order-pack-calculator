package dto

import "order-pack-calculator/internal/domain/entities"

type PackSizeResponse struct {
	ID        int64 `json:"id"`
	ProductID int   `json:"product_id"`
	Size      int   `json:"size"`
	Active    bool  `json:"active"`
}

func PackSizeResponseFromEntity(pack entities.PackSize) PackSizeResponse {
	return PackSizeResponse{
		ID:        pack.ID,
		ProductID: pack.ProductID,
		Size:      pack.Size,
		Active:    pack.Active,
	}

}

func PackSizeResponseFromEntities(packs []entities.PackSize) []PackSizeResponse {

	packSizes := make([]PackSizeResponse, 0, len(packs))

	for _, p := range packs {
		packSizes = append(packSizes, PackSizeResponseFromEntity(p))
	}
	return packSizes

}
