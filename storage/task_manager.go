package storage

import (
	_ "code-processor/docs"
	rabbitmq "code-processor/rabbitmq"
	"sync"

	"github.com/google/uuid"
)

// Task структура для задачи
// @Description Данные задачи
type Task struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Result string `json:"result,omitempty"`
}

// Управление задачами
type TaskManager struct {
	sync.Mutex
	tasks map[string]*Task
}

// экземпляр TaskManager
var TaskManagerInstance = NewTaskManager()

// фабрика таскменеджеров
func NewTaskManager() *TaskManager {
	return &TaskManager{tasks: make(map[string]*Task)}
}

// AddTask добавляет новую задачу
func (tm *TaskManager) AddTask(language, code string) string {
	taskID := uuid.New().String()
	tm.Lock()
	tm.tasks[taskID] = &Task{ID: taskID, Status: "in_progress"}
	tm.Unlock()

	/// Отправка задачи в RabbitMQ
	err := rabbitmq.SendTask(taskID, language, code)
	if err != nil {
		tm.Lock()
		tm.tasks[taskID].Status = "error"
		tm.tasks[taskID].Result = err.Error()
		tm.Unlock()
	}

	return taskID
}

// статус задачи
func (tm *TaskManager) GetTaskStatus(taskID string) (string, bool) {
	tm.Lock()
	defer tm.Unlock()
	task, exists := tm.tasks[taskID]
	if !exists {
		return "", false
	}
	return task.Status, true
}

// результат задачи
func (tm *TaskManager) GetTaskResult(taskID string) (string, bool) {
	tm.Lock()
	defer tm.Unlock()
	task, exists := tm.tasks[taskID]
	if !exists || task.Status != "ready" {
		return "", false
	}
	return task.Result, true
}

// Реализация интерфейса TaskUpdater
func (tm *TaskManager) UpdateTaskStatus(taskID, status, result string) {
	tm.Lock()
	defer tm.Unlock()
	if task, exists := tm.tasks[taskID]; exists {
		task.Status = status
		task.Result = result
	}
}
