package lock

import (
	"context"
	"errors"
)

type LockService struct {
	repo Repository
}

func NewLockService(repo Repository) *LockService {
	return &LockService{repo: repo}
}

func (s *LockService) LockCar(carID string) error {
	if carID == "" {
		return errors.New("car ID cannot be empty")
	}
	return s.repo.LockCar(context.Background(), carID)
}

func (s *LockService) UnlockCar(carID string) error {
	if carID == "" {
		return errors.New("car ID cannot be empty")
	}
	return s.repo.UnlockCar(context.Background(), carID)
}

func (s *LockService) IsCarLocked(carID string) (bool, error) {
	if carID == "" {
		return false, errors.New("car ID cannot be empty")
	}
	return s.repo.GetLockStatus(context.Background(), carID)
}
