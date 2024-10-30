package http

import (
	_ "code-processor/docs"
	"code-processor/storage"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type TaskRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// createTaskHandler создаёт новую задачу.
// @Summary Создание задачи
// @Description Создаёт новую задачу. Требуется токен аутентификации.
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {auth_token}"
// @Param task body TaskRequest true "Task Info"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /task [post]
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var taskReq TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if taskReq.Language == "" || taskReq.Code == "" {
		http.Error(w, "Language and code must be provided", http.StatusBadRequest)
		return
	}

	useTranslator(taskReq.Code)
	taskID, _ := storage.StorageInstance.TaskRepository.AddTask(taskReq.Language, taskReq.Code)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{"task_id": taskID})
}

// getTaskStatusHandler возвращает статус задачи по её ID.
// @Summary Получение статуса задачи
// @Description Получает статус задачи по ID. Требуется токен аутентификации.
// @Tags tasks
// @Security BearerAuth
// @Param Authorization header string true "Bearer {auth_token}"
// @Param task_id path string true "Task ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /status/{task_id} [get]
func getTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["task_id"]
	status, exists := storage.StorageInstance.TaskRepository.GetTaskStatus(taskID)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

// getTaskResultHandler возвращает результат задачи по её ID.
// @Summary Получение результата задачи
// @Description Получает результат задачи по ID. Требуется токен аутентификации.
// @Tags tasks
// @Security BearerAuth
// @Param Authorization header string true "Bearer {auth_token}"
// @Param task_id path string true "Task ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /result/{task_id} [get]]
func getTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["task_id"]
	result, exists := storage.StorageInstance.TaskRepository.GetTaskResult(taskID)
	if !exists {
		http.Error(w, "Task not found or not ready", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// registerUserHandler регистрирует пользователя.
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя и возвращает его ID
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRequest true "Данные для регистрации пользователя"
// @Success 201 {object} map[string]string "user_id"
// @Failure 400 {object} map[string]string "error"
// @Router /register [post]
func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req UserRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password must be provided", http.StatusBadRequest)
		return
	}

	userID, _ := storage.StorageInstance.UserRepository.AddUser(req.Login, req.Password)

	if userID == "" {
		http.Error(w, "Cannot register user", http.StatusBadRequest)
		return
	}
	response := map[string]string{
		"user_id": userID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// loginUserHandler логин пользователя.
// @Summary логин пользователя
// @Description делает авторизацию пользователя и возвращает
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRequest true "Данные для логина пользователя"
// @Success 200 {object} map[string]string "user_id"
// @Failure 400 {object} map[string]string "error"
// @Router /login [post]
func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req UserRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password must be provided", http.StatusBadRequest)
		return
	}

	userID, _ := storage.StorageInstance.UserRepository.ValidateUser(req.Login, req.Password)
	if userID == "" {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	authToken, _ := storage.StorageInstance.SessionRepository.AddSession(userID)
	if authToken == "" {
		http.Error(w, "Error with creating session", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"token": authToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
