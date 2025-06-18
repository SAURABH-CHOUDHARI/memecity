package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/SAURABH-CHOUDHARI/memecity/db/migrate"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/routes"
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

	// Initialize Redis
	redisClient := storage.NewRedisClient()

	// Set up Repository with both DB and Redis
	connections := storage.Repository{
		DB:          db,
		RedisClient: redisClient,
	}

	// Set up Fiber app
	app := fiber.New()
	app.Use(logger.New())

	// Global Rate Limiting: 100 requests per minute per IP
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return origin == "*"
		},
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,OPTIONS,DELETE", 
	}))

	// Register all routes
	routes.AllRoutes(app, connections)


	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}