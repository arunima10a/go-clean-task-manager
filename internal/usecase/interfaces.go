package usecase

import (
	"context"

	"github.com/arunima10a/task-manager/internal/entity"
)

type (
	TaskRepo interface {
		Store(ctx context.Context, t entity.Task, userID int) error
		GetAll(ctx context.Context, userID int, status string, limit int, offset int) ([]entity.Task, error)
		Update(ctx context.Context, t entity.Task, userID int) error
		Delete(ctx context.Context, id int) error
	}

	TaskUseCase interface {
		Create(ctx context.Context, t entity.Task, userID int) error
		List(ctx context.Context, userID int, status string, limit int, offset int) ([]entity.Task, error)
		UpdateStatus(ctx context.Context, id int, status string, userID int) error
		Delete(ctx context.Context, id int) error
		EnrichTaskWithQuote(ctx context.Context, taskID int) error
	}

	UserRepoInterface interface {
		Create(ctx context.Context, u entity.User) (int, error)
		GetByEmail(ctx context.Context, email string) (entity.User, error)
	}

	AuthUseCase interface {
		SignUp(ctx context.Context, email, password string) error
		Login(ctx context.Context, email, password string) (string, error)
	}

	TaskWebAPI interface {
		GetQuote(ctx context.Context) (string, error)
	}

	TransactionManager interface {
		RunInTx(ctx context.Context, f func(context.Context) error) error
	}

	TaskMessageBroker interface {
		PublishTaskCreated(taskID int, description string) error
	}
	CategoryRepo interface {
		Create(ctx context.Context, name string, userID int) (int, error)
		GetAll(ctx context.Context, userID int) ([]entity.Category, error)
	}

	CategoryUseCase interface {
		Create(ctx context.Context, name string, userID int) (int, error)
		List(ctx context.Context, userID int) ([]entity.Category, error)
	}

	TaskCache interface {
		SetTasks(ctx context.Context, userID int, tasks []entity.Task) error
		GetTasks(ctx context.Context, userID int) ([]entity.Task, error)
		DeleteTasks(ctx context.Context, userID int) error
	}
)
