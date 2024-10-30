package storage

import (
	_ "code-processor/docs"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// Session структура для хранения сессий
// @Description Данные сессии
type Session struct {
	UserID    string
	AuthToken string
	ExpiresAt time.Time
}

type SessionRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewSessionRepository(client *redis.Client, ctx context.Context) *SessionRepository {
	return &SessionRepository{client: client, ctx: ctx}
}

// AddSession создает и сохраняет сессию в Redis с токеном и сроком действия
func (sr *SessionRepository) AddSession(userID string) (string, error) {
	sessionToken := uuid.New().String()
	err := sr.client.Set(sr.ctx, sessionToken, userID, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}
	return sessionToken, nil
}

// GetUserByToken проверяет токен и возвращает ID пользователя
func (sr *SessionRepository) GetUserByToken(token string) (string, bool) {
	userID, err := sr.client.Get(sr.ctx, token).Result()
	if err == redis.Nil || err != nil {
		return "", false
	}
	return userID, true
}
