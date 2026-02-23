package v1_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"testing"

	v1 "github.com/arunima10a/task-manager/internal/controller/http/v1"
	"github.com/arunima10a/task-manager/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockTaskUseCase struct {
	listFn func(ctx context.Context, userID int, status string) ([]entity.Task, error)
}

func (m *mockTaskUseCase) Create(ctx context.Context, t entity.Task, userID int) error { return nil }
func (m *mockTaskUseCase) List(ctx context.Context, userID int, status string, limit int, offset int) ([]entity.Task, error) {
	return m.listFn(ctx, userID, status)
}
func (m *mockTaskUseCase) UpdateStatus(ctx context.Context, id int, status string, userID int) error {
	return nil
}
func (m *mockTaskUseCase) Delete(ctx context.Context, userID int) error              { return nil }
func (m *mockTaskUseCase) EnrichTaskWithQuote(ctx context.Context, taskID int) error { return nil }

func TestTaskList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUC := &mockTaskUseCase{
			listFn: func(ctx context.Context, userID int, status string) ([]entity.Task, error) {
				return []entity.Task{{Title: "Test Task"}}, nil
			},
		}
		handler := gin.New()

		handler.GET("/list", func(c *gin.Context) {
			c.Set("user_id", 1)
			c.Next()
		}, func(c *gin.Context) {
			v1.NewRouter(handler, mockUC, nil, nil, nil, nil)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/tasks/list", nil)
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Test Task")
	})
}
