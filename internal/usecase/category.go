package usecase

import (
	"context"
	"fmt"

	"github.com/arunima10a/task-manager/internal/entity"
)

type CategoryInteractor struct {
	repo CategoryRepo
}

func NewCategoryInteractor(r CategoryRepo) *CategoryInteractor {
	return &CategoryInteractor{repo: r}
}

func (uc *CategoryInteractor) Create(ctx context.Context, name string, userID int) (int, error) {
	if name == "" {
		return 0, fmt.Errorf("category name is required")
	}
	return uc.repo.Create(ctx, name, userID)
}

func (uc *CategoryInteractor) List(ctx context.Context, userID int) ([]entity.Category, error) {
	return uc.repo.GetAll(ctx, userID)
}
