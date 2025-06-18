package routes

import (
	"github.com/SAURABH-CHOUDHARI/memecity/api/handlers"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/middlewares"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/gofiber/fiber/v2"
)

func MemeRoutes(router fiber.Router, conn storage.Repository) {
	auth := router.Group("/auth", middlewares.AuthMiddleware(conn))
	memes := router.Group("/memes")

	// Public meme routes
	memes.Get("/", handlers.CreateUser(conn)) // assuming this is temp/test; rename later
	memes.Get("/leaderboard", handlers.GetLeaderboard(conn))
	memes.Get("/:id", handlers.GetMemeByID(conn))

	// Authenticated meme routes
	auth.Post("/memes", handlers.CreateMeme(conn))
	auth.Patch("/memes/:id/sale", handlers.ToggleMemeOnSale(conn))
}
