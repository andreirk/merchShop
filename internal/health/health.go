package health

import (
	"fmt"
	"log"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler(router *http.ServeMux) {
	handler := &HealthHandler{}
	router.HandleFunc("/health", handler.HealthCheck())
}

func (handler *HealthHandler) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("I'm healthy, thanks"))
		if err != nil {
			log.Println("Error while writing response in HealthCheck: %v", err)
			return
		}
		fmt.Println("Heath check")
	}
}
