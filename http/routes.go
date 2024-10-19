package http

import (
	_ "code-processor/docs"
	"code-processor/storage"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Middleware для проверки авторизации
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Разбираем токен (ожидаем Bearer {auth_token})
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]

		// Проверяем токен
		userID, valid := storage.UserManagerInstance.GetUserByToken(token)
		if !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Сохраняем userID в контексте запроса, если нужно
		r.Header.Set("User-ID", userID)

		// Передаём управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}

// маршрутизатор
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Публичные маршруты
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/register", registerUserHandler).Methods("POST")
	r.HandleFunc("/login", loginUserHandler).Methods("POST")

	// Защищённые маршруты
	protectedRoutes := r.PathPrefix("/").Subrouter()
	protectedRoutes.Use(authMiddleware)
	protectedRoutes.HandleFunc("/task", createTaskHandler).Methods("POST")
	protectedRoutes.HandleFunc("/status/{task_id}", getTaskStatusHandler).Methods("GET")
	protectedRoutes.HandleFunc("/result/{task_id}", getTaskResultHandler).Methods("GET")

	return r
}
