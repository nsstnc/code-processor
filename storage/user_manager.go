package storage

import (
	_ "code-processor/docs"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User структура для пользователя
// @Description Данные пользователя
type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
}

// Session структура для сессии
type Session struct {
	UserID    string
	AuthToken string
	ExpiresAt time.Time
}

// Управление пользователями
type UserManager struct {
	sync.Mutex
	users    map[string]*User
	sessions map[string]*Session
}

// экземпляр UserManager
var UserManagerInstance = NewUserManager()

// фабрика менеджеров
func NewUserManager() *UserManager {
	return &UserManager{
		users:    make(map[string]*User),
		sessions: make(map[string]*Session),
	}
}

// AddUser добавляет нового пользователя
func (um *UserManager) AddUser(login string, password string) string {
	um.Lock()
	defer um.Unlock()

	// Проверяем, существует ли пользователь с таким логином
	for _, user := range um.users {
		if user.Login == login {
			log.Printf("User with login '%s' already exists\n", login)
			return ""
		}
	}

	userID := uuid.New().String()

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error while hashing password: %v", err)
		return ""
	}

	// Добавляем пользователя в карту после хэширования пароля
	um.users[userID] = &User{ID: userID, Login: login, Password: string(hashedPassword)}

	return userID
}

// ValidateUser проверяет логин и пароль и создаёт сессию при успешной аутентификации
func (um *UserManager) ValidateUser(login string, password string) string {
	um.Lock()
	defer um.Unlock()

	// поиск пользователя по логину
	var foundUser *User
	for _, user := range um.users {
		if user.Login == login {
			foundUser = user
			break
		}
	}

	// Если пользователь не найден
	if foundUser == nil {
		return ""
	}

	// проверка соответствия пароля
	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	// если пароль не совпадает
	if err != nil {
		log.Printf("Invalid password for user %s\n", login)
		return ""
	}

	// Создание новой сессии и возвращение токена
	sessionToken := uuid.New().String()
	um.sessions[sessionToken] = &Session{
		UserID:    foundUser.ID,
		AuthToken: sessionToken,
		ExpiresAt: time.Now().Add(24 * time.Hour), // сессия действует 24 часа
	}

	return sessionToken
}

// GetUserByToken проверяет токен и возвращает ID пользователя
func (um *UserManager) GetUserByToken(token string) (string, bool) {
	um.Lock()
	defer um.Unlock()

	session, exists := um.sessions[token]
	if !exists || time.Now().After(session.ExpiresAt) {
		return "", false
	}

	return session.UserID, true
}
