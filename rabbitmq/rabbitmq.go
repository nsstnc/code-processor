package rabbitmq

import (
	"encoding/json"
	"log"

	"code-processor/processor"

	"github.com/streadway/amqp"
)

type Task struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

func StartConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

		result, err := processor.RunDockerContainer(task.Code, task.Language)
		if err != nil {
			log.Printf("Error running Docker container: %s", err)
			continue
		}

		log.Printf("Execution result: %s", result)
	}
}

func SendTask(language string, code string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

	task := Task{
		Language: language,
		Code:     code,
	}

	body, err := json.Marshal(task)
	if err != nil {
		log.Fatalf("Failed to marshal task: %s", err)
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
		log.Fatalf("Failed to publish a message: %s", err)
	}

	log.Printf("Sent code to RabbitMQ")
}
