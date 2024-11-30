package tokens

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
)

const SESSION_DURATION_IN_SECONDS = 3600
const SESSION_DURATION = (SESSION_DURATION_IN_SECONDS * time.Second)

func CreateSessionFromToken(token string, userID uint) error {
	sessionID := getSessionIDFromToken(token)
	dbc := db.Connect()
	session := &db.Session{
		ID:       string(sessionID),
		UserID:   userID,
		ExpireAt: time.Now().Add(SESSION_DURATION),
	}
	sessionCache.AddSession(session)
	if err := dbc.Create(session).Error; err != nil {
		return err
	}
	return nil
}

func ValidateAndGetSessionToken(token string) (bool, *db.Session) {
	sessionID := getSessionIDFromToken(token)
	dbc := db.Connect()
	session := sessionCache.GetSession(sessionID)
	if session == nil || session.ID == "" {
		log.Printf("Internal error TOKENS.SESSIONS.ValidateSessionToken.1\n")
		return false, nil
	}
	if session.ExpireAt.Unix() < time.Now().Unix() {
		dbc.Unscoped().Delete(session)
		return false, nil
	}
	return session.ID != "", session
}

var (
	ERROR_SESSION_NOT_FOUND = errors.New("Can't find session with this token")
)

func GetSession(token string) (*db.Session, error) {
	sessionID := getSessionIDFromToken(token)
	session := sessionCache.GetSession(sessionID)
	if session == nil {
		return nil, ERROR_SESSION_NOT_FOUND
	}
	return session, nil
}

func getSessionIDFromToken(token string) string {
	salt := []byte("w40DJV3v1flySvFdxHWbBSJsIOaakkVs5FG7brq4oi1#nEz2fEZxpUfyBwkkww7f")
	bytes := sha256.New().Sum(append(salt, []byte(token)...))
	return base64.URLEncoding.EncodeToString(bytes)
}
