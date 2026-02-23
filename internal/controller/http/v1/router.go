package v1

import (
	_ "github.com/arunima10a/task-manager/docs"
	"github.com/arunima10a/task-manager/internal/usecase"
	"github.com/arunima10a/task-manager/pkg/postgres"
	"github.com/arunima10a/task-manager/pkg/redis"
	"github.com/gin-gonic/gin"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func NewRouter(handler *gin.Engine, t usecase.TaskUseCase, a usecase.AuthUseCase, c usecase.CategoryUseCase, pg *postgres.Postgres, re *redis.Redis) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	newHealthRoutes(handler.Group(""), pg, re)

	h := handler.Group("/v1")

	{
		newAuthRoutes(h, a)
		taskGroup := h.Group("/tasks")
		taskGroup.Use(AuthMiddleware("your-secret-key"))

		newTaskRoutes(taskGroup, t)

	}

	categoryGroup := h.Group("/categories")

	{
		categoryGroup.Use(AuthMiddleware("your-secret-key"))
		newCategoryRoutes(categoryGroup, c)
	}

}
