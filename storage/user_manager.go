package storage

import (
	_ "code-processor/docs"
	"log"
	"sync"

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

// Управление пользователями
type UserManager struct {
	sync.Mutex
	users map[string]*User
}

// экземпляр UserManager
var UserManagerInstance = NewUserManager()

// фабрика менеджеров
func NewUserManager() *UserManager {
	return &UserManager{users: make(map[string]*User)}
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
	} else {
		um.Lock()
		um.users[userID] = &User{ID: userID, Login: login, Password: string(hashedPassword)}
		um.Unlock()
		return userID
	}
}

// ValidateUser пытается найти пользователя и проверить соответствие пароля логину
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

	// Возвращаем идентификатор пользователя
	return foundUser.ID
}
