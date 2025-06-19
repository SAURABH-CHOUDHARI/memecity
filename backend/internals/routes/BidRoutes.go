package routes

import (
	"github.com/SAURABH-CHOUDHARI/memecity/api/handlers"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/middlewares"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/gofiber/fiber/v2"
)

func BidRoutes(router fiber.Router, conn storage.Repository){
	auth:= router.Group("/auth", middlewares.AuthMiddleware(conn))
	
	auth.Post("/memes/:id/bid",handlers.PlaceBid(conn))
}