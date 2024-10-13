package http

import (
	_ "code-processor/docs"
	"code-processor/storage"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// createTaskHandler создаёт новую задачу.
// @Summary Создание задачи
// @Description Создаёт новую задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body storage.Task true "Task Info"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /task [post]
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := storage.TaskManagerInstance.AddTask()
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{"task_id": taskID})
}

// getTaskStatusHandler возвращает статус задачи по её ID.
// @Summary Получение статуса задачи
// @Description Получает статус задачи по ID
// @Tags tasks
// @Param task_id path string true "Task ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /status/{task_id} [get]
func getTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["task_id"]
	status, exists := storage.TaskManagerInstance.GetTaskStatus(taskID)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

// getTaskResultHandler возвращает результат задачи по её ID.
// @Summary Получение результата задачи
// @Description Получает результат задачи по ID
// @Tags tasks
// @Param task_id path string true "Task ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /result/{task_id} [get]]
func getTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["task_id"]
	result, exists := storage.TaskManagerInstance.GetTaskResult(taskID)
	if !exists {
		http.Error(w, "Task not found or not ready", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}
