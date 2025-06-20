package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"

	"github.com/SAURABH-CHOUDHARI/memecity/db/indexes"
	"github.com/SAURABH-CHOUDHARI/memecity/db/migrate"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/routes"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/ws"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to Postgres
	db, err := storage.NewPostgresConnection()
	if err != nil {
		log.Fatal("❌ Could not connect to DB")
	}
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	// Conditional migration
	if os.Getenv("RUN_MIGRATIONS") == "true" {
		migrate.AutoMigrate(db)
	}

	// Run index creation
	indexes.CreateIndexes(db)

	// Init Redis
	redisClient := storage.NewRedisClient()

	// Init Gemini
	ctx := context.Background()
	geminiClient, err := storage.NewGeminiClient(ctx)
	if err != nil {
		log.Fatalf("❌ Could not connect to Gemini: %v", err)
	}
	log.Println("✅ Connected to Gemini")

	// ✅ Create and start WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Final repository connections
	connections := storage.Repository{
		DB:           db,
		RedisClient:  redisClient,
		GeminiClient: geminiClient,
		Hub:          hub, // ✅ Now hub is available
	}

	// Set up Fiber
	app := fiber.New()
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return origin == "http://localhost:3000" || origin == "https://memecity-tawny.vercel.app"
		},
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,OPTIONS,DELETE,PATCH",
	}))

	// Register routes
	routes.AllRoutes(app, connections)
	app.Get("/ws", websocket.New(ws.HandleWS(hub))) // ✅ WebSocket route

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
