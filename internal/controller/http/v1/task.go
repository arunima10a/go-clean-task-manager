package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arunima10a/task-manager/internal/entity"
	"github.com/arunima10a/task-manager/internal/usecase"
	"github.com/gin-gonic/gin"
)



type taskRoutes struct {
	t usecase.TaskUseCase
}
type createTaskRequest struct {
	Title       string `json:"title"       binding:"required,min=3,max=100"  example:"Buy milk"`
	Description string `json:"description" binding:"required,max=500"  example:"At the local store"`
}

func newTaskRoutes(handler *gin.RouterGroup, t usecase.TaskUseCase) {
	r := &taskRoutes{t}

	{
		handler.GET("/list", r.list)
		handler.POST("/create", r.create)
		handler.PUT("/update/:id", r.update)
		handler.DELETE("/:id", r.delete)

	}

}

// @Summary     List tasks
// @Description get all tasks
// @Tags  	    task
// @Produce     json
// @Success     200 {array} entity.Task
// @Router      /task/list [get]
// @Security    BearerAuth
func (r *taskRoutes) list(c *gin.Context) {

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	userID := c.MustGet("user_id").(int)
	status := c.Query("status")
	tasks, err := r.t.List(c.Request.Context(), userID, status, limit, offset)
	if err != nil {

		if errors.Is(err, usecase.ErrInternal) {
			errorResponse(c, http.StatusInternalServerError, "Something went wrong on our end")
		} else {
			errorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// @Summary     Create a task
// @Description create a new task in the database
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       request body entity.Task true "Task info"
// @Success     201 {object} entity.Task
// @Router      /task/create [post]
func (r *taskRoutes) create(c *gin.Context) {
	var request createTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request body")		
		return
	}
	
	userID := c.MustGet("user_id").(int)

	task := entity.Task{
		Title: request.Title,
		Description: request.Description,
		Status: "pending",
	}
	err := r.t.Create(c.Request.Context(), task, userID)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "task service problems")		
		return
	}
	

	c.JSON(http.StatusCreated, gin.H{"message": "task created successfully"})

}
func (r *taskRoutes) update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var request struct {
		Status string `json:"status"`
		UserID int `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid Body")
		return
	}
	err := r.t.UpdateStatus(c.Request.Context(), id, request.Status, request.UserID)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task updated"})
}

func (r *taskRoutes) delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid task id")
		return
	}
	err = r.t.Delete(c.Request.Context(), id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
