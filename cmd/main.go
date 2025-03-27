package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	db "daterrr/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
)

const (
	dbSource   = "postgresql://postgres:password@172.17.0.2:5432/daterrr_db?sslmode=disable"
	serverAddr = ":4000"
)

func main() {
	// Create root context that cancels on interrupt signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	// Connect to PostgreSQL using pgx
	conn, err := pgx.Connect(ctx, dbSource)
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
	store := db.NewStore(conn)
	server := NewServer(store)

	go func() {
		log.Println("Server starting on", serverAddr)
		if err := server.router.Run(serverAddr); err != nil {
			log.Println("Server error:", err)
			cancel() 
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
}



