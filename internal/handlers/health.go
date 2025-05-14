package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HealthCheck(c *gin.Context, pool *pgxpool.Pool) {
	// Set a timeout for the database ping
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Check if the database connection is alive
	err := pool.Ping(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":   "error",
			"api_name": "Daterrr",
			"database": "unhealthy",
			"version":  "v1",
		})
		return
	}

	// If all checks pass, return healthy status
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"api_name": "Daterrr",
		"database": "healthy",
		"version":  "v1",
	})
}
