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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	store    db.Store
	router   *gin.Engine
	errorLog *log.Logger
	infoLog  *log.Logger
}

func NewServer(store *db.SQLStore, pool *pgxpool.Pool) *Server {
	gin.SetMode(gin.DebugMode)
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
	matchHandler := handlers.NewMatchHandler(store)
	chatHandler := handlers.NewChatHandler(store)
	profileHandler := handlers.NewProfileHandler(store)
	router := gin.Default()

	tMaker, err := tokengen.NewPasetoMaker(config.PasetoSecret)


	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	public := router.Group("/v1")
	{
		public.GET("/healthz", func(c *gin.Context) {
			handlers.HealthCheck(c, pool)
		})
		public.POST("/user/register", authHandler.RegisterUser)
		public.POST("/user/login", authHandler.LoginUser)
		public.GET("/chat/ws", func(c *gin.Context) {
			handlers.HandleWebSocket(c, store)
		})

	// Protected routes (require authentication)
	protected := router.Group("/v1")
	if err != nil {
		fmt.Printf("Error loading config %s", err)
	}
	protected.Use(middleware.AuthMiddleware(tMaker))
	{
		protected.GET("/genswipefeed", swipeHandler.CreateFeed)
		protected.POST("/swipes", swipeHandler.HandleSwipe)
		protected.GET("/user/getmatches", matchHandler.ListMatches)
		protected.GET("/notifications/ws", notifHandler.HandleWebSocket)
		protected.POST("/sendmessage/:receiverId", chatHandler.CreateMessage)
		protected.GET("/user/getprofile", profileHandler.GetUserProfile)
		protected.GET("/user/chats/:chatId", chatHandler.GetChatMessages)
		protected.GET("/user/getconversations", chatHandler.GetConversations)
	}

	server.router = router
	return server
}}