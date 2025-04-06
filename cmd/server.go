package main

import (
	"fmt"
	"log"
	"os"

	db "daterrr/internal/db/sqlc"
	"daterrr/internal/handlers"
	"daterrr/internal/middleware"
	"daterrr/internal/utils"
	"daterrr/pkg/auth/tokengen"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store    db.Store
	router   *gin.Engine
	errorLog *log.Logger
	infoLog  *log.Logger
}

func NewServer(store *db.SQLStore) *Server {
	config, err := utils.LoadConfig("../")
	if err != nil {
		fmt.Printf("Error loading config %s", err)
	}

	server := &Server{
		store:    store,
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(store)
	notifHandler := handlers.NewNotificationHandler(store)
	swipeHandler := handlers.NewSwipeHandler(store, notifHandler)
	router := gin.Default()
	
	public := router.Group("/v1")
	{
		public.GET("/healthcheck", handlers.HealthCheck)
		public.POST("/user/register", authHandler.RegisterUser)
		public.POST("/user/login", authHandler.LoginUser)
	}

	// Protected routes (require authentication)
	protected := router.Group("/v1")
	tMaker, err := tokengen.NewPasetoMaker(config.PasetoSecret)
	if err != nil {
		fmt.Printf("Error loading config %s", err)
	}
	protected.Use(middleware.AuthMiddleware(tMaker))
	{
		protected.POST("/swipes", swipeHandler.HandleSwipe)
		protected.GET("/notifications/ws", notifHandler.HandleWebSocket)
	}

	server.router = router
	return server
}