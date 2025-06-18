package routes

import (
	"github.com/SAURABH-CHOUDHARI/memecity/api/handlers"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/middlewares"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router, conn storage.Repository){
	auth:= router.Group("/auth", middlewares.AuthMiddleware(conn))

	router.Post("/users",handlers.CreateUser(conn))
	router.Post("/users/login", handlers.LoginUser(conn))
	auth.Post("/users/logout",handlers.LogOutUser(conn))
}