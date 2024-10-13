package http

import (
	_ "code-processor/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// маршрутизатор
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/task", createTaskHandler).Methods("POST")
	r.HandleFunc("/status/{task_id}", getTaskStatusHandler).Methods("GET")
	r.HandleFunc("/result/{task_id}", getTaskResultHandler).Methods("GET")
	return r
}
