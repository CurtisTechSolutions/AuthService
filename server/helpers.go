package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// SendJSONResponseHelper is a helper function to send JSON responses.
//
// It sets the Content-Type header to application/json and writes the response.
// It also sets the status code based on the success field of the response.
func SendJSONResponseHelper(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")

	if response.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		slog.Error("Error marshalling JSON response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
