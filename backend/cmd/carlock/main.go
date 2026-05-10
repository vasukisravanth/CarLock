package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"car-lock-system/backend/internal/api"
	"car-lock-system/backend/internal/config"
	"car-lock-system/backend/internal/db"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Initialize MongoDB connection
	mongoURI := cfg.MongoDB.URI
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Default to localhost
	}

	err = db.InitMongoDB(mongoURI, cfg.MongoDB.DBName)
	if err != nil {
		log.Fatalf("could not connect to MongoDB: %v", err)
	}

	// Initialize repositories
	api.InitRepositories()

	// Set up router
	router := mux.NewRouter()
	api.RegisterRoutes(router)

	// Server address
	serverAddr := ":" + cfg.Port
	if cfg.Port == "" {
		serverAddr = ":8080" // Default port
	}

	// Start server
	log.Printf("Starting server on %s...", serverAddr)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	// Graceful shutdown
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := db.CloseMongoConnection(ctx); err != nil {
			log.Printf("error closing MongoDB connection: %v", err)
		}
	}()
}
