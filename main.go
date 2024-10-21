package main

import (
	routes "code-processor/http"
	rabbitmq "code-processor/rabbitmq"
	"code-processor/storage"
	"log"
	"net/http"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		rabbitmq.StartConsumer(storage.TaskManagerInstance)
	}()

	log.Println("Waiting for consumer to start...")

	// Запускаем HTTP-сервер
	r := routes.NewRouter()
	log.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))

	// Ждем завершения потребителя перед выходом
	wg.Wait()
}
