package session

import "database/sql"

type Session struct {
	Id        int
	SessionId string
	UserId    int
}

type Model struct {
	db *sql.DB
}

func NewSessionModel(db *sql.DB) Model {
	return Model{db}
}

func (sm *Model) Add(session Session) (Session, error) {
	result, err := sm.db.Exec("INSERT INTO sessions (session_id, user_id) VALUES (?, ?, ?)", session.SessionId, session.UserId)

	if err != nil {
		return Session{}, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return Session{}, err
	}

	session.Id = int(id)

	return session, nil
}

func (sm *Model) IsValid(sessionId string) bool {
	if err := sm.db.QueryRow("SELECT id FROM sessions WHERE session_id = ?", sessionId); err != nil {
		return false
	}

	return true
}
