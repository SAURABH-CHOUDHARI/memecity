package services

import (
	"errors"
	"os"
	"time"
	"context"
	"fmt"
	"encoding/json"
	"strings"

	"github.com/SAURABH-CHOUDHARI/memecity/internals/models"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(conn storage.Repository, email, username, password string) (string, error) {
	if email == "" || username == "" || password == "" {
		return "", errors.New("email, username, and password are required")
	}

	// Check if user already exists
	var existing models.User
	if err := conn.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return "", errors.New("user with this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	// 🔐 Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	newUser := models.User{
		ID:       uuid.New(),
		Email:    email,
		Username: username,
		Password: string(hashedPassword), // ✅ Save hashed password
	}

	if err := conn.DB.Create(&newUser).Error; err != nil {
		return "", err
	}

	// Generate JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	claims := jwt.MapClaims{
		"sub":   newUser.ID.String(),
		"email": newUser.Email,
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
		"jti":   uuid.NewString(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	// Cache user in Redis
	userKey := fmt.Sprintf("user:%s", newUser.ID.String())
	userJSON, err := json.Marshal(newUser)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user: %w", err)
	}

	if err := conn.RedisClient.Set(
		context.Background(),
		userKey,
		userJSON,
		3*24*time.Hour,
	).Err(); err != nil {
		return "", fmt.Errorf("failed to store user in redis: %w", err)
	}

	return signedToken, nil
}

func LogoutUser(conn storage.Repository, authHeader string) error {
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return errors.New("missing or invalid authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return errors.New("JWT_SECRET not set")
	}

	// Parse JWT to get claims
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return errors.New("missing jti in token")
	}

	expUnix, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("missing exp in token")
	}
	ttl := time.Until(time.Unix(int64(expUnix), 0))

	// Blacklist the token
	ctx := context.Background()
	blacklistKey := fmt.Sprintf("blacklist:%s", jti)

	if err := conn.RedisClient.Set(ctx, blacklistKey, "true", ttl).Err(); err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

func LoginUser(conn storage.Repository, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email and password are required")
	}

	// Find user by email
	var user models.User
	if err := conn.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	// 🔐 Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Generate JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	jti := uuid.NewString()
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"jti":   jti,
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	// Cache user in Redis
	userKey := fmt.Sprintf("user:%s", user.ID.String())
	userJSON, _ := json.Marshal(user)
	if err := conn.RedisClient.Set(context.Background(), userKey, userJSON, 3*24*time.Hour).Err(); err != nil {
		return "", fmt.Errorf("failed to store user in redis: %w", err)
	}

	return signedToken, nil
}