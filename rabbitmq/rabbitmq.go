package rabbitmq

import (
	"code-processor/processor"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/streadway/amqp"
)

type Task struct {
	ID       string `json:"id"`
	Language string `json:"language"`
	Code     string `json:"code"`
}

type TaskUpdater interface {
	UpdateTaskStatus(taskID string, status string, result string)
}

func StartConsumer(updater TaskUpdater) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"code_tasks",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)

		var task Task
		if err := json.Unmarshal(msg.Body, &task); err != nil {
			log.Printf("Error unmarshalling task: %s", err)
			continue
		}
		task.Code = strings.ReplaceAll(task.Code, "\\n", "\n")
		// Обработка задачи
		result, err := processor.RunDockerContainer(task.Code, task.Language)

		// Обновление статуса через переданный updater
		if err != nil {
			updater.UpdateTaskStatus(task.ID, "error", err.Error())
		} else {
			updater.UpdateTaskStatus(task.ID, "ready", result)
		}

		log.Printf("Execution result for task %s: %s", task.ID, result)
	}
}

func SendTask(taskID string, language string, code string) error {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"code_tasks",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	task := Task{
		ID:       taskID,
		Language: language,
		Code:     code,
	}

	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf("Sent task %s to RabbitMQ", taskID)
	return nil
}
