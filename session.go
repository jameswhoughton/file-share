package file_share

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateKey() string {
	key := make([]byte, 32)
	rand.Read(key)

	return base64.StdEncoding.EncodeToString(key)
}

type Session struct {
	Id        int
	SessionId string
	UserId    int
}

type SessionService interface {
	Add(session Session) (Session, error)
	IsValid(sessionId string) bool
}
