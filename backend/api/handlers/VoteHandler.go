package handlers

import (
	"errors"

	"github.com/SAURABH-CHOUDHARI/memecity/internals/models"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/services"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type voteRequest struct {
	Type string `json:"type"` // "up" or "down"
}

func VoteMeme(conn storage.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		memeIDParam := c.Params("id")
		memeID, err := uuid.Parse(memeIDParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid meme ID",
			})
		}

		var req voteRequest
		if err := c.BodyParser(&req); err != nil || (req.Type != "up" && req.Type != "down") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid vote type, must be 'up' or 'down'",
			})
		}

		user := c.Locals("user").(models.User)

		if err := services.VoteOnMeme(conn, user.ID, memeID, req.Type); err != nil {
			code := fiber.StatusInternalServerError
			if errors.Is(err, services.ErrInvalidVoteType) {
				code = fiber.StatusBadRequest
			} else if errors.Is(err, services.ErrMemeNotFound) {
				code = fiber.StatusNotFound
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "vote recorded",
		})
	}
}

