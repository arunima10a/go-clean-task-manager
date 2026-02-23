package v1

import (
	"net/http"

	"github.com/arunima10a/task-manager/internal/usecase"
	"github.com/gin-gonic/gin"
)

type categoryRoutes struct {
	uc usecase.CategoryUseCase
}

func newCategoryRoutes(handler *gin.RouterGroup, uc usecase.CategoryUseCase) {
	r := &categoryRoutes{uc}
	{
		handler.POST("", r.create)
		handler.GET("", r.list)
	}
}

func (r *categoryRoutes) create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	userID := c.MustGet("user_id").(int)
	id, err := r.uc.Create(c.Request.Context(), req.Name, userID)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to create category")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (r *categoryRoutes) list(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	categories, err := r.uc.List(c.Request.Context(), userID)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to fetch categories")
		return
	}
	c.JSON(http.StatusOK, categories)
}
