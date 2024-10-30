package storage

import (
	_ "code-processor/docs"
	rabbitmq "code-processor/rabbitmq"
	"database/sql"
	"log"
	"strings"

	"github.com/google/uuid"
)

// Task структура для задачи
// @Description Данные задачи
type Task struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Result string `json:"result,omitempty"`
}

type TaskRepository struct {
	db *sql.DB
}

// Инициализация TaskRepository с подключением к базе данных
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// AddTask добавляет новую задачу в PostgreSQL и возвращает её ID
func (tr *TaskRepository) AddTask(language, code string) (string, error) {
	taskID := uuid.New().String()

	// Вставка задачи в базу данных
	_, err := tr.db.Exec("INSERT INTO tasks (id, status) VALUES ($1, $2)", taskID, "in_progress")
	if err != nil {
		log.Println("Failed to add task:", err)
		return "", err
	}

	/// Отправка задачи в RabbitMQ
	err = rabbitmq.SendTask(taskID, language, code)
	if err != nil {
		// Обновляем статус задачи на "error" в базе данных в случае ошибки
		_, updateErr := tr.db.Exec("UPDATE tasks SET status = $1, result = $2 WHERE id = $3", "error", err.Error(), taskID)
		if updateErr != nil {
			log.Println("Failed to update task status:", updateErr)
			return "", updateErr
		}
		return "", err
	}

	return taskID, nil
}

// GetTaskStatus возвращает статус задачи по ID
func (tr *TaskRepository) GetTaskStatus(taskID string) (string, bool) {
	var status string
	err := tr.db.QueryRow("SELECT status FROM tasks WHERE id = $1", taskID).Scan(&status)
	if err == sql.ErrNoRows {
		return "", false
	} else if err != nil {
		log.Println("Failed to get task status:", err)
		return "", false
	}
	return status, true
}

// GetTaskResult возвращает результат задачи по ID, если она завершена
func (tr *TaskRepository) GetTaskResult(taskID string) (string, bool) {
	var result string
	var status string
	err := tr.db.QueryRow("SELECT status, result FROM tasks WHERE id = $1", taskID).Scan(&status, &result)
	if err == sql.ErrNoRows || status != "ready" {
		return "", false
	} else if err != nil {
		log.Println("Failed to get task result:", err)
		return "", false
	}
	return result, true
}

// UpdateTaskStatus обновляет статус и результат задачи в PostgreSQL
func (tr *TaskRepository) UpdateTaskStatus(taskID, status, result string) error {
	result = sanitizeString(result)

	_, err := tr.db.Exec("UPDATE tasks SET status = $1, result = $2 WHERE id = $3", status, result, taskID)
	if err != nil {
		log.Println("Failed to update task status:", err)
		return err
	}
	return nil
}

func sanitizeString(s string) string {
	return strings.ReplaceAll(s, "\x00", "") // Удаляем нулевые байты
}
