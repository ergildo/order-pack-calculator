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

// Constructor for PackSizeService
func NewPackSizeService(packSizeRepository repositories.PackSizeRepository) PackSizeService {
	return packSizeService{packSizeRepository: packSizeRepository}
}

type packSizeService struct {
	packSizeRepository repositories.PackSizeRepository
}

// Creates a new pack size entry
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

// Updates an existing pack size
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

// Retrieves all pack sizes
func (p packSizeService) GetAll(ctx context.Context) ([]dto.PackSizeResponse, error) {
	packSizes, err := p.packSizeRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch pack size. %w", err)
	}
	responses := dto.PackSizeResponseFromEntities(packSizes)
	return responses, nil
}

// Calculate optimal pack sizes for an order
func (p packSizeService) CalcOptimalPacks(ctx context.Context, order dto.CalculatePackSizesRequest) (*dto.OptimalPackSizesResponse, error) {
	packSizes, err := p.packSizeRepository.GetSizesByProductID(ctx, int64(order.ProductID))
	if err != nil {
		return nil, fmt.Errorf("could not fetch pack sizes. %w", err)
	}

	solution := p.calcOptimalPacks(order, packSizes)
	return solution, nil
}

// Core logic: calculates optimal pack combination using dynamic programming
func (packSizeService) calcOptimalPacks(order dto.CalculatePackSizesRequest, packSizes []int) *dto.OptimalPackSizesResponse {
	maxPackSize := slices.Max(packSizes)
	limit := order.OrderQuantity + maxPackSize

	// dp[i] holds the best solution to exactly pack i items
	dp := make([]*dto.OptimalPackSizesResponse, limit+1)
	dp[0] = &dto.OptimalPackSizesResponse{} // base case: 0 items needs 0 packs

	for i := 0; i <= limit; i++ {
		if dp[i] == nil {
			continue
		}
		for _, size := range packSizes {
			next := i + size
			if next > limit {
				continue
			}

			// Clone current pack combination and add one more pack of current size
			newCombo := cloneCombination(dp[i].PackCombination)
			addToCombination(&newCombo, size)

			newSolution := &dto.OptimalPackSizesResponse{
				PackCombination: newCombo,
				TotalItems:      next,
				TotalPacks:      dp[i].TotalPacks + 1,
			}

			// Update dp[next] if it's a better solution (fewer items or fewer packs)
			if dp[next] == nil ||
				newSolution.TotalItems < dp[next].TotalItems ||
				(newSolution.TotalItems == dp[next].TotalItems && newSolution.TotalPacks < dp[next].TotalPacks) {
				dp[next] = newSolution
			}
		}
	}

	// Find best valid solution with minimum total items >= order quantity
	best := &dto.OptimalPackSizesResponse{TotalItems: math.MaxInt64}
	for i := order.OrderQuantity; i <= limit; i++ {
		if dp[i] != nil && (dp[i].TotalItems < best.TotalItems ||
			(dp[i].TotalItems == best.TotalItems && dp[i].TotalPacks < best.TotalPacks)) {
			best = dp[i]
		}
	}

	return best
}

// Deep copies a pack combination slice
func cloneCombination(combo []dto.PackDetail) []dto.PackDetail {
	var cloned []dto.PackDetail
	for _, cb := range combo {
		cloned = append(cloned, dto.PackDetail{Size: cb.Size, Count: cb.Count})
	}
	return cloned
}

// Adds a pack of the given size to the combination (increments if already exists)
func addToCombination(combo *[]dto.PackDetail, size int) {
	for i := range *combo {
		if (*combo)[i].Size == size {
			(*combo)[i].Count++
			return
		}
	}
	*combo = append(*combo, dto.PackDetail{Size: size, Count: 1})
}
