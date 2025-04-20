package db

import (
	"fmt"
	"time"

	"github.com/CTS/AuthService/internal"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UserID    uint
	SessionID string
	ExpiresAt time.Time
}

func createSessionID(email string) string {
	// Create a unique session ID based on the email and current time
	sessionID := internal.EncodeSHA256([]byte(email + time.Now().String()))
	return fmt.Sprintf("%x", sessionID)
}

func SessionCreate(user *User, expiresAt time.Duration) (string, error) {
	sessionID := createSessionID(user.Email)
	expiresAtTime := time.Now().Add(expiresAt)
	result := DB.Create(&Session{
		UserID:    user.ID,
		SessionID: sessionID,
		ExpiresAt: expiresAtTime,
	})
	if result.Error != nil {
		return sessionID, result.Error
	}

	return sessionID, nil
}

func SessionGet(sessionID string) (*Session, error) {
	var session Session
	result := DB.Where("session_id = ?", sessionID).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func SessionValidate(sessionID string) (bool, error) {
	var session Session
	result := DB.Where("session_id = ? AND expires_at > ?", sessionID, time.Now()).First(&session)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func SessionExpire(sessionID string) error {
	var session Session
	result := DB.Where("session_id = ?", sessionID).First(&session)
	if result.Error != nil {
		return result.Error
	}

	session.ExpiresAt = time.Now()
	result = DB.Save(&session)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
