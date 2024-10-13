package http

import (
	"github.com/gorilla/mux"
)

// маршрутизатор
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/task", createTaskHandler).Methods("POST")
	r.HandleFunc("/status/{task_id}", getTaskStatusHandler).Methods("GET")
	r.HandleFunc("/result/{task_id}", getTaskResultHandler).Methods("GET")
	return r
}
