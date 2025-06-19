package handlers

import (
	"errors"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/models"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/services"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/gofiber/fiber/v2"
)

type PlaceBidRequest struct {
	Credits int `json:"credits"`
}

func PlaceBid(conn storage.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req PlaceBidRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		user := c.Locals("user").(models.User)
		memeID := c.Params("id")
		if memeID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid meme ID",
			})
		}

		if err := services.PlaceBid(conn, user.ID, memeID, req.Credits); err != nil {
			code := fiber.StatusInternalServerError
			if errors.Is(err, services.ErrInvalidBid) || errors.Is(err, services.ErrInsufficientCredits) {
				code = fiber.StatusBadRequest
			} else if errors.Is(err, services.ErrNotFound) || errors.Is(err, services.ErrUnauthorized) {
				code = fiber.StatusForbidden
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"memeId": memeID, 
			"message": "success",
		})
	}
}
