package http

import (
	"code-processor/storage"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// новая задача
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := storage.TaskManagerInstance.AddTask()
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"id": taskID})
}

// статус задачи
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

// результат задачи
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
