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

func CreateUser(conn storage.Repository, email, username, password string) (string, models.User, error) {
	if email == "" || username == "" || password == "" {
		return "", models.User{}, errors.New("email, username, and password are required")
	}

	var existing models.User
	if err := conn.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return "", models.User{}, errors.New("user with this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", models.User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", models.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := models.User{
		ID:       uuid.New(),
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}

	if err := conn.DB.Create(&newUser).Error; err != nil {
		return "", models.User{}, err
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", models.User{}, errors.New("JWT_SECRET not set")
	}

	claims := jwt.MapClaims{
		"sub":   newUser.ID.String(),
		"email": newUser.Email,
		"jti":   uuid.NewString(),
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", models.User{}, err
	}

	userJSON, _ := json.Marshal(newUser)
	conn.RedisClient.Set(context.Background(), fmt.Sprintf("user:%s", newUser.ID.String()), userJSON, 3*24*time.Hour)

	return signedToken, newUser, nil
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

func LoginUser(conn storage.Repository, email, password string) (string, models.User, error) {
	if email == "" || password == "" {
		return "", models.User{}, errors.New("email and password are required")
	}

	var user models.User
	if err := conn.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", models.User{}, errors.New("user not found")
		}
		return "", models.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", models.User{}, errors.New("invalid password")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", models.User{}, errors.New("JWT_SECRET not set")
	}

	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"jti":   uuid.NewString(),
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", models.User{}, err
	}

	userKey := fmt.Sprintf("user:%s", user.ID.String())
	userJSON, _ := json.Marshal(user)
	conn.RedisClient.Set(context.Background(), userKey, userJSON, 3*24*time.Hour)

	return signedToken, user, nil
}
