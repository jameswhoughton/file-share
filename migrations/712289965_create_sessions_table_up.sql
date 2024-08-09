CREATE TABLE sessions (
	id INTEGER PRIMARY KEY,
	session_id TEXT NOT NULL,
    user_id INTEGER,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id)
)