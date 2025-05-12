package db

import (
	"time"

	"github.com/CTS/AuthService/internal"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UserID    uint
	SessionID string
	ExpiresAt time.Time
	Expired   bool
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

// expireSession invalidates a session in the DB.
// It does this by setting the expiration time of the session to the current time, and sets expired to true.
func (s *Session) expireSession() {
	s.ExpiresAt = time.Now()
	s.Expired = true
}

func createSessionID(email string) string {
	// Create a unique session ID based on the email and current time
	return internal.EncodeSHA256([]byte(email + time.Now().String()))
}

func SessionCreate(user *User, expiresAt time.Duration) (string, error) {
	// Create a new session
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
	result := DB.Where(&Session{SessionID: sessionID, Expired: false}).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func SessionValidate(sessionID string) (bool, error) {
	var session Session
	// result := DB.Where("session_id = ? AND expires_at > ?", sessionID, time.Now()).First(&session)
	result := DB.Where(&Session{
		SessionID: sessionID,
		Expired:   false,
	}).First(&session)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func SessionExpire(sessionID string) error {
	var session Session
	result := DB.Where(&Session{SessionID: sessionID, Expired: false}).First(&session)
	if result.Error != nil {
		return result.Error
	}

	session.expireSession()
	result = DB.Save(&session)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// SessionInvalidateForUser invalidates all sessions for a given user
// by setting their expiration time to the current time.
func SessionExpireAllByUserID(userID uint) error {
	var sessions []Session
	result := DB.Where(&Session{
		UserID:  userID,
		Expired: false,
	}).Find(&sessions)
	if result.Error != nil {
		return result.Error
	}

	for _, session := range sessions {
		session.expireSession()
		result = DB.Save(&session)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// SessionExpireAllByUserID invalidates all sessions for a given user ID.
// We do this by getting the user ID from the session ID and then calling SessionInvalidateForUser.
// This will invalidate all sessions for the user including the current session.
func SessionExpireAllBySessionID(sessionID string) error {
	var session Session
	result := DB.Where("session_id = ?", sessionID).First(&session)
	if result.Error != nil {
		return result.Error
	}

	err := SessionExpireAllByUserID(session.UserID)
	if err != nil {
		return err
	}

	return nil
}
