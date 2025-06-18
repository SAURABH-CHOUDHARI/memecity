package routes

import (
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func AllRoutes(app *fiber.App, conn storage.Repository) {
    // Base API group
	api := app.Group("/api")
    
    
    UserRoutes(api, conn)
	MemeRoutes(api, conn)

}