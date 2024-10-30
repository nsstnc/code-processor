package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

type Storage struct {
	UserRepository    *UserRepository
	TaskRepository    *TaskRepository
	SessionRepository *SessionRepository
}

// Глобальная переменная StorageInstance, представляющая единственный экземпляр Storage
var StorageInstance *Storage

func CreateTables(db *sql.DB) error {
	userTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY,
        login VARCHAR(50) UNIQUE NOT NULL,
        password VARCHAR(100) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	taskTableSQL := `
    CREATE TABLE IF NOT EXISTS tasks (
        id UUID PRIMARY KEY,
        status TEXT,
		result TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := db.Exec(userTableSQL)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	_, err = db.Exec(taskTableSQL)
	if err != nil {
		return fmt.Errorf("error creating tasks table: %w", err)
	}

	return nil
}

// InitStorage инициализирует StorageInstance с подключением к PostgreSQL и Redis
func InitStorage() {
	config, err := LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
	)

	// Подключение к PostgreSQL
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Проверка подключения к PostgreSQL
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	// Инициализация Redis
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Redis.Addr,
	})

	// Проверка подключения к Redis
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	if err := CreateTables(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Создание экземпляра Storage с инициализированными репозиториями
	StorageInstance = &Storage{
		UserRepository:    NewUserRepository(db),
		TaskRepository:    NewTaskRepository(db),
		SessionRepository: NewSessionRepository(redisClient, ctx),
	}

	log.Println("Storage successfully initialized")
}
