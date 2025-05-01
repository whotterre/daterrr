package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time" 

	db "daterrr/internal/db/sqlc"
	"daterrr/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
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

	poolConfig, err := pgxpool.ParseConfig(config.DBSource)
	if err != nil {
		log.Fatal("Couldn't parse database config:", err)
	}

	// --- Configure Pool Settings (Adjust as needed) ---
	poolConfig.MaxConns = int32(25)         // Max connections in pool
	poolConfig.MinConns = int32(5)          // Min connections (optional)
	poolConfig.MaxConnLifetime = 5 * time.Minute // Connections older than this are recycled
	poolConfig.MaxConnIdleTime = 1 * time.Minute // Connections idle longer than this are closed
	// --- End Pool Settings ---


	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatal("Couldn't create database connection pool:", err)
	}
	defer pool.Close()


	// Verify connection pool
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Database pool connection failed:", err)
	}
	log.Print("Successfully connected to the PostgreSQL database pool!")

	store := db.NewStore(pool).(*db.SQLStore) // Pass the POOL
	

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
