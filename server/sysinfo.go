package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func InfoRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/hostname", systemHostname)
	return r
}

func systemHostname(w http.ResponseWriter, r *http.Request) {
	// Get the system hostname
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("Unable to get hostname", "error", err.Error())
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Unable to get hostname",
		})
		return
	}
	// Write the hostname to the response
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Hostname retrieved successfully",
		Data:    hostname,
	})
}
