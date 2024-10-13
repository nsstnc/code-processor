package storage

import (
	_ "code-processor/docs"
	"math/rand"
	"sync"
	"time"

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
func (tm *TaskManager) AddTask() string {
	taskID := uuid.New().String()
	tm.Lock()
	tm.tasks[taskID] = &Task{ID: taskID, Status: "in_progress"}
	tm.Unlock()

	// обработка
	go func() {
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second) // случайная задержка от 1 до 5 секунд
		tm.Lock()
		tm.tasks[taskID].Status = "ready"
		tm.tasks[taskID].Result = "Some result data" // Здесь будет результат обработки
		tm.Unlock()
	}()

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
