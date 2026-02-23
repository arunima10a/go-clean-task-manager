package v1

import (
	"context"
	"errors"

	pb "github.com/arunima10a/task-manager/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/arunima10a/task-manager/internal/entity"
	"github.com/arunima10a/task-manager/internal/usecase"
)

type TaskHandler struct {
	pb.UnimplementedTaskServiceServer

	t usecase.TaskUseCase
	c usecase.CategoryUseCase
}

func NewTaskHandler(t usecase.TaskUseCase, c usecase.CategoryUseCase) *TaskHandler {
	return &TaskHandler{
		t: t,
		c: c,
	}
}

func (h *TaskHandler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not found in context")
	}

	task := entity.Task{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
	}

	err := h.t.Create(ctx, task, userID)
	if err != nil {
		return nil, parseError(err)
	}

	return &pb.CreateTaskResponse{
		Status: "Task created successfully via gRPC",
	}, nil
}
func (h *TaskHandler) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	tasks, err := h.t.List(ctx, 1, req.GetStatus(), 10, 0)
	if err != nil {
		return nil, parseError(err)
	}

	var grpcTasks []*pb.Task
	for _, t := range tasks {
		grpcTasks = append(grpcTasks, &pb.Task{
			Id:           int32(t.ID),
			Title:        t.Title,
			Description:  t.Description,
			Status:       string(t.Status),
			CategoryId:   int32(t.CategoryID),
			CategoryName: t.CategoryName,
		})

	}
	return &pb.ListTasksResponse{
		Tasks: grpcTasks,
	}, nil
}

func (h *TaskHandler) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	id, err := h.c.Create(ctx, req.GetName(), 1)
	if err != nil {
		return nil, err
	}
	return &pb.CreateCategoryResponse{Id: int32(id)}, nil
}

func (h *TaskHandler) ListCategories(ctx context.Context, req *pb.ListCategoryRequest) (*pb.ListCategoryResponse, error) {
	cats, err := h.c.List(ctx, 1)
	if err != nil {
		return nil, err
	}
	var grpcCats []*pb.Category
	for _, val := range cats {
		grpcCats = append(grpcCats, &pb.Category{
			Id:   int32(val.ID),
			Name: val.Name,
		})
	}
	return &pb.ListCategoryResponse{Categories: grpcCats}, nil

}

func parseError(err error) error {
	if errors.Is(err, usecase.ErrTaskNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	return status.Error(codes.Internal, "internal system error")
}
