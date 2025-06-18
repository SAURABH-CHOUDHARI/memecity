package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

// ---- Redis Init ----

func NewRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_URL") // should be just host:port
	if addr == "" {
		addr = "localhost:6379"
	}

	pass := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: "default", 
		Password: pass,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("❌ Failed to connect to Redis at %s: %v", addr, err)
	} else {
		log.Printf("✅ Connected to Redis at %s", addr)
	}

	return client
}

// ---- Helper ----

func getEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return v
}

// ---- DB Init ----

func NewPostgresConnection() (*gorm.DB, error) {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(getEnvInt("DB_MAX_IDLE_CONNS", 50))
	sqlDB.SetMaxOpenConns(getEnvInt("DB_MAX_OPEN_CONNS", 500))
	sqlDB.SetConnMaxLifetime(time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME", 10)) * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return db, nil
}

// ---- Final Init Function ----

func NewRepository() (*Repository, error) {
	db, err := NewPostgresConnection()
	if err != nil {
		return nil, err
	}

	redisClient := NewRedisClient()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Repository{
		DB:          db,
		RedisClient: redisClient,
	}, nil
}
