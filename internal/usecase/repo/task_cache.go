package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arunima10a/task-manager/internal/entity"
	"github.com/arunima10a/task-manager/pkg/redis"
)

type TaskCacheRepo struct {
	redis *redis.Redis
}

func NewTaskCache(r *redis.Redis) *TaskCacheRepo {
	return &TaskCacheRepo{redis: r}
}

func (r *TaskCacheRepo) GetTasks(ctx context.Context, userID int) ([]entity.Task, error) {
	key := fmt.Sprintf("tasks:%d", userID)
	val, err := r.redis.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var tasks []entity.Task
	err = json.Unmarshal([]byte(val), &tasks)
	return tasks, err
}

func (r *TaskCacheRepo) SetTasks(ctx context.Context, userID int, tasks []entity.Task) error {
	key := fmt.Sprintf("tasks:%d", userID)
	data, _ := json.Marshal(tasks)
	// Cache for 5 minutes
	return r.redis.Client.Set(ctx, key, data, 5*time.Minute).Err()
}

func (r *TaskCacheRepo) DeleteTasks(ctx context.Context, userID int) error {
	key := fmt.Sprintf("tasks:%d", userID)
	return r.redis.Client.Del(ctx, key).Err()
}