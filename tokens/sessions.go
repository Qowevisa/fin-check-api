package tokens

import (
	"crypto/sha256"
	"log"
	"time"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
)

func getSalt() []byte {
	return []byte("w40DJV3v1flySvFdxHWbBSJsIOaakkVs5FG7brq4oi1#nEz2fEZxpUfyBwkkww7f")
}

func CreateSessionFromToken(token string, userID uint) error {
	salt := getSalt()
	sessionID := sha256.New().Sum(append(salt, []byte(token)...))
	dbc := db.Connect()
	session := &db.Session{
		ID:       string(sessionID),
		UserID:   userID,
		ExpireAt: time.Now().Add(time.Hour),
	}
	if err := dbc.Create(session).Error; err != nil {
		return err
	}
	return nil
}

func ValidateSessionToken(token string) bool {
	salt := getSalt()
	sessionID := sha256.New().Sum(append(salt, []byte(token)...))
	dbc := db.Connect()
	session := &db.Session{}
	if err := dbc.Find(session, sessionID).Error; err != nil {
		log.Printf("DBERROR: %v\n", err)
		return false
	}
	if session.ID == "" {
		return false
	}
	if session.ExpireAt.Unix() < time.Now().Unix() {
		dbc.Delete(session)
		return false
	}
	return session.ID != ""
}

func GetSession(token string) (*db.Session, error) {
	salt := getSalt()
	sessionID := sha256.New().Sum(append(salt, []byte(token)...))
	dbc := db.Connect()
	session := &db.Session{}
	if err := dbc.Find(session, sessionID).Error; err != nil {
		return nil, err
	}
	return session, nil
}
