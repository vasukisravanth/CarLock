package lock

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	LockCar(ctx context.Context, carID string) error
	UnlockCar(ctx context.Context, carID string) error
	GetLockStatus(ctx context.Context, carID string) (bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) LockCar(ctx context.Context, carID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE cars SET locked = TRUE WHERE id = $1", carID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UnlockCar(ctx context.Context, carID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE cars SET locked = FALSE WHERE id = $1", carID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetLockStatus(ctx context.Context, carID string) (bool, error) {
	var locked bool
	err := r.db.QueryRowContext(ctx, "SELECT locked FROM cars WHERE id = $1", carID).Scan(&locked)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return locked, nil
}
