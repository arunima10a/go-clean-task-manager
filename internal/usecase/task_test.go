package usecase_test

import (
	"context"
	"sync"
	"testing"

	"github.com/arunima10a/task-manager/internal/entity"
	"github.com/arunima10a/task-manager/internal/usecase"
	"github.com/arunima10a/task-manager/pkg/logger"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	tasks []entity.Task
}

func (m *mockRepo) Store(ctx context.Context, t entity.Task, userID int) error {
	t.UserID = userID
	m.tasks = append(m.tasks, t)
	return nil
}

func (m *mockRepo) GetAll(ctx context.Context, userID int, status string, limit int, offset int) ([]entity.Task, error) {
	return m.tasks, nil
}

func (m *mockRepo) Update(ctx context.Context, t entity.Task, userID int) error {
	return nil
}

func (m *mockRepo) Delete(ctx context.Context, id int) error {
	return nil
}

type mockWebAPI struct{}

func (m *mockWebAPI) GetQuote(ctx context.Context) (string, error) {
	return "Test Quote", nil
}

type mockBroker struct{}

func (m *mockBroker) PublishTaskCreated(id int, desc string) error { return nil }

type mockCache struct{}
 func(m *mockCache) SetTasks(ctx context.Context, userID int, tasks []entity.Task) error {return nil}
 func(m *mockCache) GetTasks(ctx context.Context, userID int) ([]entity.Task, error) {return nil, nil}
 func(m *mockCache) DeleteTasks(ctx context.Context, userID int) error {return nil}

func TestCreateTask(t *testing.T) {
	repo := &mockRepo{}
	l := logger.New("error")
	web := &mockWebAPI{}
	var wg sync.WaitGroup
	broker := &mockBroker{}
	cache := &mockCache{}

	uc := usecase.New(repo, web, l, &wg, broker, cache)
	mockUserID := 1

	t.Run("success", func(t *testing.T) {
		task := entity.Task{Title: "Test Task", Status: "active"}
		err := uc.Create(context.Background(), task, mockUserID)

		assert.NoError(t, err)
		assert.Len(t, repo.tasks, 1)
		assert.Equal(t, "Test Task", repo.tasks[0].Title)
		assert.Equal(t, mockUserID, repo.tasks[0].UserID)

	})

	t.Run("empty title error", func(t *testing.T) {
		task := entity.Task{Title: "", Status: "active"}
		err := uc.Create(context.Background(), task, mockUserID)

		assert.Error(t, err)
	})
}
