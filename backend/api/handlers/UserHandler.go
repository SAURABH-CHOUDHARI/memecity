package handlers

import (
	"github.com/SAURABH-CHOUDHARI/memecity/internals/models"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/services"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/gofiber/fiber/v2"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(conn storage.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateUserRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		token, user, err := services.CreateUser(conn, req.Email, req.Username, req.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":         user.ID,
				"email":      user.Email,
				"username":   user.Username,
				"profilePic": user.ProfilePic,
				"credits":    user.Credits,
			},
		})
	}
}


func LogOutUser(conn storage.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if err := services.LogoutUser(conn, authHeader); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		user := c.Locals("user").(models.User)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "logged out",
			"email":    user.Email,
			"username": user.Username,
		})
	}
}

func LoginUser(conn storage.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req loginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		token, user, err := services.LoginUser(conn, req.Email, req.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":         user.ID,
				"email":      user.Email,
				"username":   user.Username,
				"profilePic": user.ProfilePic,
				"credits":    user.Credits,
			},
		})
	}
}
