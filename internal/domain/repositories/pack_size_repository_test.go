package repositories

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"order-pack-calculator/internal/domain/entities"
	errs "order-pack-calculator/internal/domain/errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPackSizeRepository(db)

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO pack_sizes (product_id, size)
VALUES ($1, $2)
RETURNING id, active`)).
			WithArgs(int64(1), 10).
			WillReturnRows(sqlmock.NewRows([]string{"id", "active"}).AddRow(100, true))

		res, err := repo.Create(context.Background(), entities.PackSize{ProductID: 1, Size: 10})
		assert.NoError(t, err)
		assert.Equal(t, int64(100), res.ID)
		assert.True(t, res.Active)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO pack_sizes (product_id, size) VALUES ($1, $2) RETURNING id, active")).
			WithArgs(int64(2), 20).
			WillReturnError(errors.New("insert error"))

		_, err := repo.Create(context.Background(), entities.PackSize{ProductID: 2, Size: 20})
		assert.Error(t, err)
	})
}

func TestUpdate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPackSizeRepository(db)

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE pack_sizes SET size = $1, active = $2 WHERE id = $3")).
			WithArgs(20, true, int64(1)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(context.Background(), entities.PackSize{ID: 1, Size: 20, Active: true})
		assert.NoError(t, err)
	})

	t.Run("no rows affected", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE pack_sizes SET size = $1, active = $2 WHERE id = $3")).
			WithArgs(15, false, int64(99)).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := repo.Update(context.Background(), entities.PackSize{ID: 99, Size: 15, Active: false})
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("exec error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE pack_sizes SET size = $1, active = $2 WHERE id = $3")).
			WithArgs(10, true, int64(2)).
			WillReturnError(errors.New("update error"))

		err := repo.Update(context.Background(), entities.PackSize{ID: 2, Size: 10, Active: true})
		assert.Error(t, err)
	})
}

func TestGetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPackSizeRepository(db)

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, product_id, size, active FROM pack_sizes WHERE id = $1")).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "size", "active"}).AddRow(1, 1, 10, true))

		res, err := repo.GetByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, product_id, size, active FROM pack_sizes WHERE id = $1")).
			WithArgs(int64(2)).
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetByID(context.Background(), 2)
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})
}

func TestGetSizesByProductID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPackSizeRepository(db)

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT size FROM pack_sizes WHERE product_id = $1 AND active = true")).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"size"}).AddRow(10).AddRow(20))

		sizes, err := repo.GetSizesByProductID(context.Background(), 1)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []int{10, 20}, sizes)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT size FROM pack_sizes WHERE product_id = $1 AND active = true")).
			WithArgs(int64(2)).
			WillReturnError(errors.New("query failed"))

		_, err := repo.GetSizesByProductID(context.Background(), 2)
		assert.Error(t, err)
	})
}
