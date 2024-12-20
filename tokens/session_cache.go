package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"strings"
	"sync"
	"time"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
)

type Token struct {
	Id         uint
	Val        string
	LastActive time.Time
}

type SessiomMapMu struct {
	Initialized bool
	SessionMap  map[string]*db.Session
	Mu          sync.RWMutex
}

var sessionCache SessiomMapMu

// NOTE: should be launch with a goroutine
// NOTE: it cannot die
func Init() error {
	sessionCache.SessionMap = make(map[string]*db.Session)
	var dbSessions []*db.Session
	dbc := db.Connect()
	if err := dbc.Find(&dbSessions).Error; err != nil {
		return err
	}
	timeNow := time.Now()
	for _, dbSession := range dbSessions {
		if dbSession.ExpireAt.Unix() < timeNow.Unix() {
			if err := dbc.Unscoped().Delete(dbSession).Error; err != nil {
				log.Printf("dbc.Unscoped().Delete(dbSession) error: %v\n", err)
			}
			continue
		}
		sessionCache.SessionMap[dbSession.ID] = dbSession
	}
	sessionCache.Initialized = true
	return nil
}

func (s *SessiomMapMu) HaveSession(sessionID string) bool {
	s.Mu.RLock()
	_, exists := s.SessionMap[sessionID]
	s.Mu.RUnlock()
	return exists
}

func (s *SessiomMapMu) AddSession(session *db.Session) {
	s.Mu.Lock()
	s.SessionMap[session.ID] = session
	s.Mu.Unlock()
}

func (s *SessiomMapMu) GetSession(sessionID string) *db.Session {
	s.Mu.RLock()
	val, exists := s.SessionMap[sessionID]
	s.Mu.RUnlock()
	if !exists {
		return nil
	}
	return val
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("generateRandomString ERROR: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func generateTokenVal() string {
	for {
		tok := generateRandomString(32)
		trimedToken := strings.Trim(tok, "=")
		// TODO: do some thing so it can check if user will have the same token
		return trimedToken
	}
}

func AddToken(id uint) (*Token, error) {
	val := generateTokenVal()
	token := &Token{
		Id:         id,
		Val:        val,
		LastActive: time.Now(),
	}
	return token, nil
}
