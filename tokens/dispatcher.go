package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"sync"
	"time"
)

type Token struct {
	Val        string
	LastActive time.Time
}

var (
	ActiveDur = time.Duration(time.Hour)
)

func (t Token) IsExpired() bool {
	return time.Now().Sub(t.LastActive) >= ActiveDur
}

type TokensMapMu struct {
	Initialized bool
	Tokmap      map[uint]Token
	TokmapRev   map[string]Token
	Mu          sync.RWMutex
}

var toks TokensMapMu

// NOTE: should be launch with a goroutine
// NOTE: it cannot die
func StartTokens() {
	if toks.Initialized {
		return
	}
	toks.Tokmap = make(map[uint]Token)
	toks.TokmapRev = make(map[string]Token)
	toks.Initialized = true
	for {
		//
		toks.Mu.Lock()
		for id, token := range toks.Tokmap {
			if token.IsExpired() {
				delete(toks.Tokmap, id)
			}
		}
		toks.Mu.Unlock()
		//
		time.Sleep(time.Minute)
	}
}

var (
	ERROR_DONT_HAVE_TOKEN    = errors.New("Don't have token for this user")
	ERROR_ALREADY_HAVE_TOKEN = errors.New("Already have token")
)

func GetToken(id uint) (Token, error) {
	toks.Mu.RLock()
	val, exists := toks.Tokmap[id]
	toks.Mu.RUnlock()
	if !exists {
		return Token{}, ERROR_DONT_HAVE_TOKEN
	}
	return val, nil
}

func haveToken(id uint) bool {
	toks.Mu.RLock()
	_, exists := toks.Tokmap[id]
	toks.Mu.RUnlock()
	return exists
}

func UpdateLastActive(id uint) error {
	if !haveToken(id) {
		return ERROR_DONT_HAVE_TOKEN
	}
	toks.Mu.Lock()
	val := toks.Tokmap[id]
	val.LastActive = time.Now()
	toks.Tokmap[id] = val
	toks.Mu.Unlock()
	return nil
}

func haveTokenVal(val string) bool {
	toks.Mu.RLock()
	_, exists := toks.TokmapRev[val]
	toks.Mu.RUnlock()
	return exists
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("generateRandomString: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func generateTokenVal() string {
	for {
		tok := generateRandomString(32)
		if !haveTokenVal(tok) {
			return tok
		}
	}
}

func AddToken(id uint) (Token, error) {
	toks.Mu.RLock()
	_, exists := toks.Tokmap[id]
	toks.Mu.RUnlock()
	if exists {
		return Token{}, ERROR_ALREADY_HAVE_TOKEN
	}
	return Token{
		Val:        generateTokenVal(),
		LastActive: time.Now(),
	}, nil
}
