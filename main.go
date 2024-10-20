package main

import (
	routes "code-processor/http"
	rabbitmq "code-processor/rabbitmq"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	language := "cpp"
	code := `#include <iostream>

		int main() {
			std::cout << "Hello, World!" << std::endl;
			return 0;
		}`

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		rabbitmq.StartConsumer()
	}()

	log.Println("Waiting for consumer to start...")
	// Добавляем задержку, чтобы потребитель успел подключиться
	time.Sleep(time.Second)

	// Отправляем задачу
	rabbitmq.SendTask(language, code)

	// Запускаем HTTP-сервер
	r := routes.NewRouter()
	log.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))

	// Ждем завершения потребителя перед выходом
	wg.Wait()
}
