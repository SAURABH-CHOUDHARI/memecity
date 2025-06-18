package routes

import (
	"github.com/SAURABH-CHOUDHARI/memecity/api/handlers"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/middlewares"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/gofiber/fiber/v2"
)

func MemeRoutes(router fiber.Router, conn storage.Repository){
	auth:= router.Group("/auth", middlewares.AuthMiddleware(conn))

	router.Get("/memes/:id", handlers.GetMemeByID(conn))
	router.Get("/memes",handlers.CreateUser(conn))
	auth.Post("/memes",handlers.CreateMeme(conn))
	auth.Patch("/memes/:id/sale", handlers.ToggleMemeOnSale(conn))
}