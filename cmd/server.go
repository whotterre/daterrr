package main

import (
	db "daterrr/internal/db/sqlc"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store    db.Store
	router   *gin.Engine
	errorLog *log.Logger
	infoLog  *log.Logger
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	// Define routes here
	router.GET("/v1/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"api_name": "Daterrr",
			"version": "v1",
		})
	})
	server.router = router
	return server
}
