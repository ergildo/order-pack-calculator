package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"order-pack-calculator/internal/domain/entities"
	errs "order-pack-calculator/internal/domain/errors"
)

func NewPackSizeRepository(db *sql.DB) PackSizeRepository {
	return packSizeRepository{db: db}
}

type packSizeRepository struct {
	db *sql.DB
}


func (p packSizeRepository) Create(ctx context.Context, pack entities.PackSize) (*entities.PackSize, error) {
	query := `
	INSERT INTO pack_sizes (product_id, size)
	VALUES ($1, $2)
	RETURNING id, active
`

	err := p.db.QueryRowContext(ctx, query, pack.ProductID, pack.Size).Scan(&pack.ID, &pack.Active)
	if err != nil {
		return nil, fmt.Errorf("failed to insert pack size for product_id=%d, size=%d: %w", pack.ProductID, pack.Size, err)
	}
	return &pack, nil
}
func (p packSizeRepository) Update(ctx context.Context, pack entities.PackSize) error {
	query := `
		UPDATE pack_sizes
		SET size = $1, active = $2
		WHERE id = $3
	`
	rs, err := p.db.ExecContext(ctx, query, pack.Size, pack.Active, pack.ID)
	if err != nil {
		return fmt.Errorf("failed to update pack size id=%d: %w", pack.ID, err)
	}
	rowsAffected, err := rs.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update pack size id=%d: %w", pack.ID, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: id=%d", errs.ErrNotFound, pack.ID)
	}

	return nil
}
func (p packSizeRepository) GetSizesByProductID(ctx context.Context, productID int64) ([]int, error) {
	query := `
	SELECT size
	FROM pack_sizes
	WHERE product_id = $1 AND active = true
`
	rows, err := p.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query pack sizes. product_id=%d: %w", productID, err)
	}
	defer rows.Close()

	var packSizes []int
	for rows.Next() {
		var size int
		if err := rows.Scan(&size); err != nil {
			return nil, fmt.Errorf("failed to scan pack size row: %w", err)
		}
		packSizes = append(packSizes, size)
	}

	return packSizes, nil
}

func (p packSizeRepository) GetByID(ctx context.Context, ID int64) (*entities.PackSize, error) {
	query := `
	SELECT id, product_id, size, active
	FROM pack_sizes
	WHERE id = $1
`
	var packSize entities.PackSize
	err := p.db.QueryRowContext(ctx, query, ID).Scan(&packSize.ID, &packSize.ProductID, &packSize.Size, &packSize.Active)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errs.ErrNotFound
		default:
			return nil, err
		}
	}

	return &packSize, nil
}


// GetAll implements PackSizeRepository.
func (p packSizeRepository) GetAll(ctx context.Context) ([]entities.PackSize, error) {
	query := `
	SELECT id, product_id, size, active
	FROM pack_sizes
`
rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pack sizes. %w", err)
	}
	defer rows.Close()

	var packSizes []entities.PackSize
	for rows.Next() {
		var packSize entities.PackSize
		if err := rows.Scan(&packSize.ID, &packSize.ProductID, &packSize.Size, &packSize.Active); err != nil {
			return nil, fmt.Errorf("failed to scan pack size row: %w", err)
		}
		packSizes = append(packSizes, packSize)
	}

	return packSizes, nil
}
