package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	file_share "github.com/jameswhoughton/file-share"
)

type SessionService struct {
	db *sql.DB
}

func NewSessionService(db *sql.DB) SessionService {
	return SessionService{db}
}

func (ss *SessionService) Add(session file_share.Session) (file_share.Session, error) {
	result, err := ss.db.Exec("INSERT INTO sessions (session_id, user_id) VALUES (?, ?)", session.SessionId, session.UserId)

	if err != nil {
		return file_share.Session{}, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return file_share.Session{}, err
	}

	session.Id = int(id)

	return session, nil
}

func (ss *SessionService) Destroy(sessionId string) error {
	_, err := ss.db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionId)

	if err != nil {
		return fmt.Errorf("could not destroy session: %v", err)
	}

	return nil
}

func (ss *SessionService) IsValid(sessionId string) bool {
	row := ss.db.QueryRow("SELECT id FROM sessions WHERE session_id = ?", sessionId)

	if err := row.Scan(); errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}
