package storage

import (
	_ "code-processor/docs"
	"database/sql"
	"log"

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

type UserRepository struct {
	db *sql.DB
}

// Инициализация UserRepository с подключением к базе данных
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// AddUser добавляет нового пользователя в PostgreSQL
func (ur *UserRepository) AddUser(login string, password string) (string, error) {
	// Проверка существования пользователя с таким логином
	var userExists bool
	err := ur.db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE login = $1)", login).Scan(&userExists)
	if err != nil {
		return "", err
	}
	if userExists {
		log.Printf("User with login '%s' already exists\n", login)
		return "", nil
	}

	userID := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error while hashing password: %v", err)
		return "", err
	}

	// Вставка пользователя в базу данных
	_, err = ur.db.Exec("INSERT INTO users (id, login, password) VALUES ($1, $2, $3)", userID, login, hashedPassword)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// ValidateUser проверяет логин и пароль
func (ur *UserRepository) ValidateUser(login string, password string) (string, bool) {
	var user User
	err := ur.db.QueryRow("SELECT id, password FROM users WHERE login = $1", login).Scan(&user.ID, &user.Password)
	if err != nil {
		return "", false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Invalid password for user %s\n", login)
		return "", false
	}

	return user.ID, true
}
