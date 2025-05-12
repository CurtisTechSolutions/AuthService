package internal

import (
	"crypto/sha256"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(plainTextPassword string) (string, error) {
	if plainTextPassword == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(hashedPassword string, password string) bool {
	if hashedPassword == "" || password == "" {
		slog.Error("hashed password or password is empty")
		return false
	}

	// Compare the hashed password with the provided password
	// err is nil if the passwords match
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		slog.Error("failed to verify password", slog.String("error", err.Error()))
	}

	return err == nil
}

func EncodeSHA256(data []byte) string {
	if len(data) == 0 {
		slog.Error("data is empty")
		return ""
	}

	// Compute the SHA256 hash
	hash := sha256.Sum256(data)
	// the hash in hexadecimal format
	return fmt.Sprintf("%x", hash)
}
