package main

import (
	routes "code-processor/http"
	"log"
	"net/http"
)

func main() {
	r := routes.NewRouter()
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
