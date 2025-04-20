package server

import (
	"log/slog"
	"net/http"
	"os"
)

func systemHostname(w http.ResponseWriter, r *http.Request) {
	// Get the system hostname
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("Unable to get hostname", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Write the hostname to the response
	w.Write([]byte(hostname))
	w.WriteHeader(http.StatusOK)
}
