package main

import (
	"log"
	"net/http"
	"sync"

	routes "code-processor/http"
	"code-processor/rabbitmq"
	"code-processor/storage"
)

func main() {
	storage.InitStorage()

	// Запуск RabbitMQ consumer в отдельной горутине
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		rabbitmq.StartConsumer(storage.StorageInstance.TaskRepository) // Передаем taskRepo в качестве TaskUpdater
	}()

	log.Println("Waiting for consumer to start...")

	// Запускаем HTTP-сервер с созданием маршрутов и передачей репозиториев
	r := routes.NewRouter()
	r.Handle("/metrics", routes.PrometheusHandler())

	log.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))

	// Ждем завершения потребителя перед выходом
	wg.Wait()
}
