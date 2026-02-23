package v1

import (
	"net/http"

	"github.com/arunima10a/task-manager/pkg/postgres"
	"github.com/arunima10a/task-manager/pkg/redis"
	"github.com/gin-gonic/gin"
)

func newHealthRoutes(handler *gin.RouterGroup, pg *postgres.Postgres, re *redis.Redis) {
	handler.GET("/health", func(c *gin.Context) {
		if err := pg.Pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "database": "error"})
			return
		}

		if err := re.Client.Ping(c.Request.Context()).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "redis": "error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})
}