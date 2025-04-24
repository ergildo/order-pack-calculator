package services

import (
	"context"
	"fmt"
	"math"
	"order-pack-calculator/internal/domain/dto"
	"order-pack-calculator/internal/domain/entities"
	"slices"

	"order-pack-calculator/internal/domain/repositories"
)

func NewPackSizeService(packSizeRepository repositories.PackSizeRepository) PackSizeService {
	return packSizeService{packSizeRepository: packSizeRepository}
}

type packSizeService struct {
	packSizeRepository repositories.PackSizeRepository
}

func (p packSizeService) Create(ctx context.Context, request dto.CreatePackSizeRequest) (*dto.PackSizeResponse, error) {

	packSize := entities.PackSize{
		ProductID: request.ProductID,
		Size:      request.Size,
	}

	saved, err := p.packSizeRepository.Create(ctx, packSize)
	if err != nil {
		return nil, fmt.Errorf("could not update pack size. %w", err)
	}

	response := dto.PackSizeResponseFromEntity(*saved)

	return &response, nil
}

func (p packSizeService) Update(ctx context.Context, request dto.UpdatePackSizeRequest) error {
	packSize, err := p.packSizeRepository.GetByID(ctx, request.ID)
	if err != nil {
		return fmt.Errorf("could not update pack size. %w", err)
	}

	if request.Size != nil {
		packSize.Size = *request.Size
	}
	if request.Active != nil {
		packSize.Active = *request.Active
	}

	err = p.packSizeRepository.Update(ctx, *packSize)

	if err != nil {
		return fmt.Errorf("could not update pack size. %w", err)
	}

	return nil
}

func (p packSizeService) GetAll(ctx context.Context) ([]dto.PackSizeResponse, error) {
	packSizes, err := p.packSizeRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch pack size. %w", err)
	}
	responses := dto.PackSizeResponseFromEntities(packSizes)
	return responses, nil
}

func (p packSizeService) CalcOptimalPacks(ctx context.Context, order dto.CalculatePackSizesRequest) (*dto.OptimalPackSizesResponse, error) {
	packSizes, err := p.packSizeRepository.GetSizesByProductID(ctx, int64(order.ProductID))
	if err != nil {
		return nil, fmt.Errorf("could not fetch pack sizes. %w", err)
	}

	solution := p.calcOptimalPacks(order, packSizes)
	return solution, nil
}

func (packSizeService) calcOptimalPacks(order dto.CalculatePackSizesRequest, packSizes []int) *dto.OptimalPackSizesResponse {

	maxPackSize := slices.Max(packSizes)

	limit := order.OrderQuantity + maxPackSize

	// dp[i] = best solution for exactly i items
	dp := make([]*dto.OptimalPackSizesResponse, limit+1)
	dp[0] = &dto.OptimalPackSizesResponse{}

	for i := 0; i <= limit; i++ {
		if dp[i] == nil {
			continue
		}
		for _, size := range packSizes {
			next := i + size
			if next > limit {
				continue
			}

			// clone current combination
			newCombo := cloneCombination(dp[i].PackCombination)
			addToCombination(&newCombo, size)

			newSolution := &dto.OptimalPackSizesResponse{
				PackCombination: newCombo,
				TotalItems:      i + size,
				TotalPacks:      dp[i].TotalPacks + 1,
			}

			if dp[next] == nil ||
				newSolution.TotalItems < dp[next].TotalItems ||
				(newSolution.TotalItems == dp[next].TotalItems && newSolution.TotalPacks < dp[next].TotalPacks) {
				dp[next] = newSolution
			}
		}
	}

	// find the best solution with minimal totalItems >= orderQuantity
	best := &dto.OptimalPackSizesResponse{TotalItems: math.MaxInt64}
	for i := order.OrderQuantity; i <= limit; i++ {
		if dp[i] != nil && (dp[i].TotalItems < best.TotalItems ||
			(dp[i].TotalItems == best.TotalItems && dp[i].TotalPacks < best.TotalPacks)) {
			best = dp[i]
		}
	}

	return best

}

func cloneCombination(combo []dto.PackDetail) []dto.PackDetail {
	var cloned []dto.PackDetail
	for _, cb := range combo {
		cloned = append(cloned, dto.PackDetail{Size: cb.Size, Count: cb.Count})
	}
	return cloned
}

func addToCombination(combo *[]dto.PackDetail, size int) {
	for i := range *combo {
		if (*combo)[i].Size == size {
			(*combo)[i].Count++
			return
		}
	}
	*combo = append(*combo, dto.PackDetail{Size: size, Count: 1})
}
