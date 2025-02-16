package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Error while encoding json: %v", err)
		return
	}
}
