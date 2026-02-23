package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/arunima10a/task-manager/internal/entity"
	"github.com/arunima10a/task-manager/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var tasksCreated = promauto.NewCounter(prometheus.CounterOpts{
	Name: "task_manager_tasks_created_total",
	Help: "The total number of created tasks",
})

type TaskInteractor struct {
	repo   TaskRepo
	web    TaskWebAPI
	l      *logger.Logger
	wg     *sync.WaitGroup
	broker TaskMessageBroker
	cache  TaskCache
}

func New(r TaskRepo, w TaskWebAPI, l *logger.Logger, wg *sync.WaitGroup, b TaskMessageBroker, c TaskCache) *TaskInteractor {
	return &TaskInteractor{
		repo:   r,
		web:    w,
		l:      l,
		wg:     wg,
		broker: b,
		cache:  c,
	}
}

func (uc *TaskInteractor) Create(ctx context.Context, t entity.Task, userID int) error {

	if t.Title == "" {
		return fmt.Errorf("TaskUseCase - Create: title is required")
	}
	quote, err := uc.web.GetQuote(ctx)
	if err != nil {
		uc.l.Error(err, "TaskUseCase - Create - uc.web.GetQuote")
		quote = "Keep moving forward!"
	}
	t.Description = fmt.Sprintf("%s\n\nDaily Motivation: %s", t.Description, quote)

	err = uc.repo.Store(ctx, t, userID)
	if err != nil {
		uc.l.Error(err, "TaskUseCase - Create - uc.repo.Store")
		return fmt.Errorf("TaskUseCase - Create - uc.repo.Store: %w", err)
	}
	tasksCreated.Inc()

	return uc.broker.PublishTaskCreated(t.ID, t.Description)
}

func (uc *TaskInteractor) List(ctx context.Context, userID int, status string, limit int, offset int) ([]entity.Task, error) {
	cachedTasks, err := uc.cache.GetTasks(ctx, userID)
	if err == nil && len(cachedTasks) > 0 {
		uc.l.Info("Cache Hit for user %d", userID)
		return cachedTasks, nil
	}

	uc.l.Info("Cache Miss for user %d, hitting Postgres", userID)
	tasks, err := uc.repo.GetAll(ctx, userID, status, limit, offset)
	if err != nil {
		return nil, err
	}

	_ = uc.cache.SetTasks(ctx, userID, tasks)

	return tasks, nil
}
func (uc *TaskInteractor) UpdateStatus(ctx context.Context, id int, status string, userID int) error {
	if id <= 0 {

		return fmt.Errorf("invalid task id")
	}
	task := entity.Task{
		ID:     id,
		Status: status,
	}
	return uc.repo.Update(ctx, task, userID)
}

func (uc *TaskInteractor) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("TaskUseCase - Delete: invalid id")
	}
	return uc.repo.Delete(ctx, id)
}
func (uc *TaskInteractor) EnrichTaskWithQuote(ctx context.Context, taskID int) error {

	uc.wg.Add(1)
	defer uc.wg.Done()
	quote, err := uc.web.GetQuote(ctx)
	if err != nil {
		return fmt.Errorf("TaskUseCase - EnrichTaskWithQueue - GetQuote: %w", err)
	}
	task := entity.Task{
		ID:          taskID,
		Description: "Quote: " + quote,
	}
	err = uc.repo.Update(ctx, task, 0)
	if err != nil {
		uc.l.Error(err, "TaskUseCase - EnrichTaskWithQuote - Update")
		return err
	}

	uc.l.Info("Task %d enriched with quote via RabbitMQ Worker", taskID)
	return nil
}
