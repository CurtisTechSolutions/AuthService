package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/CTS/AuthService/db"
	"github.com/CTS/AuthService/internal"
	"github.com/go-chi/chi/v5"
)

func AuthRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/login", login)
	r.Post("/signup", signup)
	r.Post("/logout", logout)
	r.Post("/validate", validateSession)
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
	if err := BodyParser(r, &form); err != nil {
		slog.Error("Error parsing request body", "error", err.Error())
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Invalid request",
		})
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
	if err := BodyParser(r, &form); err != nil {
		slog.Error("Error parsing request body", "error", err.Error())
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	// Check if email is already in use
	exists, err := db.UserExists(&db.User{Email: form.Email})
	if err != nil {
		slog.Error("Error checking user existence", "error", err.Error(), "email", form.Email)
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Unable to create user",
		})
		return
	}
	if exists {
		slog.Error("Email already in use", "email", form.Email)
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Email already in use",
		})
		return
	}

	// Validate the password
	if form.Password == "" {
		slog.Error("Password cannot be empty")
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Please check your email and password",
		})
		return
	}

	// Check if the password is strong enough
	if len(form.Password) < 8 {
		slog.Error("Password must be at least 8 characters long")
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Password must be at least 8 characters long",
		})
		return
	}

	// Check if the email is valid
	if form.Email == "" {
		slog.Error("Email cannot be empty")
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Please check your email and password",
		})
		return
	}
	// Hash the password before storing it
	hashedPassword, err := internal.HashPassword(form.Password)
	if err != nil {
		slog.Error("Error hashing password", "error", err.Error())
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Unable to create user",
		})
		return
	}
	// Check if the hashed password is valid
	if hashedPassword == "" {
		slog.Error("Error hashing password")
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Unable to create user",
		})
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
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Unable to create user",
		})
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	// Remove any session cookies
	slog.Debug("Implement me :)")
	w.Write([]byte("Logged out"))

	// Find all sessions for the user in the database and delete them or set them to expired
}

func validateSession(w http.ResponseWriter, r *http.Request) {
	// Check if the session cookie is present
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		slog.Error("Session cookie not found", "error", err.Error())
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Session cookie not found",
		})
		return
	}

	// Validate the session in the database
	session, err := db.SessionValidate(sessionCookie.Value)
	if err != nil && !session {
		slog.Error("Error validating session", "error", err.Error())
		SendJSONResponse(w, Response{
			Success: false,
			Message: "Session not found",
		})
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Message: "Session is valid",
	})
}
