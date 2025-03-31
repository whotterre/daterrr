package main

import (
	db "daterrr/internal/db/sqlc"
	"daterrr/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store    db.Store
	router   *gin.Engine
	errorLog *log.Logger
	infoLog  *log.Logger
}

func NewServer(store *db.SQLStore) *Server {
	server := &Server{store: store}
	authHandler := handlers.NewAuthHandler(store)
	router := gin.Default()
	// Define routes here
	router.GET("/v1/healthcheck", handlers.HealthCheck)
	router.POST("/v1/user/register", authHandler.RegisterUser)
	router.POST("/v1/user/login", authHandler.LoginUser)
	server.router = router
	return server
}
