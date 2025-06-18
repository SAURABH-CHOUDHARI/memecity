package handlers

import (
	"github.com/SAURABH-CHOUDHARI/memecity/internals/models"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/services"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type createMemeRequest struct {
	Title    string   `json:"title"`
	ImageURL string   `json:"image_url"`
	Tags     []string `json:"tags"`
	Caption  string   `json:"caption"`
}

func CreateMeme(conn storage.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req createMemeRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		owner, ok := user.(models.User)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "user context error",
			})
		}

		err := services.CreateMeme(conn, req.Title, req.ImageURL, req.Caption, req.Tags, owner.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "meme created",
		})
	}
}

func ToggleMemeOnSale(conn storage.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		memeIDStr := c.Params("id") // expects /meme/:id/sale
		memeID, err := uuid.Parse(memeIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid meme ID",
			})
		}

		var meme models.Meme
		if err := conn.DB.First(&meme, "id = ?", memeID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "meme not found",
			})
		}

		user := c.Locals("user").(models.User)
		if meme.OwnerID != user.ID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "not the owner of this meme",
			})
		}

		if meme.OnSale {
			if err := conn.DB.Model(&meme).Update("on_sale", false).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to update sale status",
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "meme is now off sale",
			})
		}

		// Toggle OnSale to true
		if err := conn.DB.Model(&meme).Update("on_sale", true).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to update sale status",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "meme is now on sale",
			})
	}
}

func GetMemeByID(conn storage.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		memeIDStr := c.Params("id")
		memeID, err := uuid.Parse(memeIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid meme ID",
			})
		}

		var meme models.Meme
		if err := conn.DB.Preload("Owner").First(&meme, "id = ?", memeID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "meme not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(meme)
	}
}