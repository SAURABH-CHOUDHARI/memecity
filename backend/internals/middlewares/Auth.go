package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SAURABH-CHOUDHARI/memecity/internals/models"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware(conn storage.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing or invalid authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "JWT_SECRET not set",
			})
		}

		// Parse and validate JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["sub"] == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		//  Check if token is blacklisted
		jti, ok := claims["jti"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing jti in token",
			})
		}

		blacklistKey := fmt.Sprintf("blacklist:%s", jti)
		ctx := context.Background()
		isBlacklisted, err := conn.RedisClient.Get(ctx, blacklistKey).Result()
		if err == nil && isBlacklisted == "true" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token is blacklisted",
			})
		}

		// ✅ Extract user ID
		userID, ok := claims["sub"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid user ID in token",
			})
		}

		// Try Redis for user data
		userKey := fmt.Sprintf("user:%s", userID)
		userJSON, err := conn.RedisClient.Get(ctx, userKey).Result()
		var user models.User

		if err == nil {
			if err := json.Unmarshal([]byte(userJSON), &user); err == nil {
				c.Locals("user", user)
				return c.Next()
			}
		}

		// Redis miss → fallback to DB
		if err := conn.DB.First(&user, "id = ?", userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "user not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "db error: " + err.Error(),
			})
		}

		// Re-cache in Redis
		userBytes, _ := json.Marshal(user)
		conn.RedisClient.Set(ctx, userKey, userBytes, 7*24*time.Hour)

		c.Locals("user", user)
		return c.Next()
	}
}

