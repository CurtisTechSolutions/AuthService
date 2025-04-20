package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/CTS/AuthService/db"
	"github.com/CTS/AuthService/internal"
	"github.com/go-chi/chi/v5"
)

func BaseRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/login", login)
	r.Post("/signup", signup)
	r.Post("/logout", logout)
	r.Post("/validate", validateSession)
	r.Get("/dbg/hostname", systemHostname)
	return r
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
	// Implementation for logging in a user
	// Get username and password from request
	// Validate the username and password
	var form LoginForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the user exists in the database
	user, err := db.UserGet(&db.User{Email: form.Email, Status: "active"})
	if err != nil {
		slog.Error("Error getting user", "error", err.Error(), "email", form.Email)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Check if user exists and is not active
	if user == nil {
		slog.Error("User not found", "email", form.Email)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the user entered the correct password
	if !internal.VerifyPassword(user.Password, form.Password) {
		slog.Error("Invalid password", "email", form.Email)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a session and add it to the database. Expires in 24 hours
	sessionID, err := db.SessionCreate(user, time.Hour*24)
	if err != nil {
		slog.Error("Error creating session", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Send the session back to the user
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	})
}

type SignupForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func signup(w http.ResponseWriter, r *http.Request) {
	// Implementation for signing up a user
	var form SignupForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// Check if email is already in use
	exists, err := db.UserExists(&db.User{Email: form.Email})
	if err != nil {
		slog.Error("Error checking user existence", "error", err.Error(), "email", form.Email)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if exists {
		slog.Error("Email already in use", "email", form.Email)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Email already in use"))
		return
	}

	// Validate the password
	if form.Password == "" {
		slog.Error("Password cannot be empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Check if the password is strong enough
	if len(form.Password) < 8 {
		slog.Error("Password must be at least 8 characters long")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Check if the email is valid
	if form.Email == "" {
		slog.Error("Email cannot be empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Hash the password before storing it
	hashedPassword, err := internal.HashPassword(form.Password)
	if err != nil {
		slog.Error("Error hashing password", "error", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Check if the hashed password is valid
	if hashedPassword == "" {
		slog.Error("Error hashing password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create a new user in the database
	err = db.UserCreate(&db.User{
		Status:   "active",
		Email:    form.Email,
		Password: hashedPassword,
		Role:     "user",
	})
	if err != nil {
		slog.Error("Error creating user", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	// Remove any session cookies
	slog.Debug("Implement me :)")
	w.Write([]byte("Logged out"))

	// Find all sessions for the user in the database and delete them or set them to expired
}

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

func validateSession(w http.ResponseWriter, r *http.Request) {
	// Check if the session cookie is present
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		slog.Error("Session cookie not found", "error", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Validate the session in the database
	session, err := db.SessionValidate(sessionCookie.Value)
	if err != nil && !session {
		slog.Error("Error validating session", "error", err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Session not found"))
		return
	}

	w.Write([]byte("Session is valid"))
}
