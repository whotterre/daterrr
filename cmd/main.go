package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	db "daterrr/internal/db/sqlc"
	"daterrr/internal/utils"

	"github.com/jackc/pgx/v5"
)

func main() {
	configPath := filepath.Join("../")
	config, err := utils.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Error loading config file", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	// Connect to PostgreSQL using pgx
	conn, err := pgx.Connect(ctx, config.DBSource)
	if err != nil {
		log.Fatal("Couldn't connect to database:", err)
	}
	defer conn.Close(ctx)
	
	// Verify connection
	if err := conn.Ping(ctx); err != nil {
		log.Fatal("Database connection failed:", err)
	}
	log.Print("Successfully connected to the PostgreSQL database!")
	// Create store and server
	store := db.NewStore(conn).(*db.SQLStore) 
	server := NewServer(store)

	go func() {
		log.Println("Server starting on", config.ServerAddr)
		if err := server.router.Run(config.ServerAddr); err != nil {
			log.Println("Server error:", err)
			cancel() 
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
}



