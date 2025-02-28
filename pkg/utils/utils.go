package utils

import (
	"encoding/json"
	"go-flashcards-server/pkg/types"
	"net/http"
)

func HandleErrorResponse(w http.ResponseWriter, message string, code int) {
	response := types.GCResponse[string]{
		IsOK:    false,
		Message: message,
		Payload: nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
