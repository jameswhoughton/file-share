package sqlite

import (
	"database/sql"
	"fmt"

	file_share "github.com/jameswhoughton/file-share"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) UserService {
	return UserService{db}
}

func (us *UserService) Get(id int) (file_share.User, error) {
	var user file_share.User

	if err := us.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user); err != nil {
		if err == sql.ErrNoRows {
			return file_share.User{}, fmt.Errorf("user %d not found", id)
		}

		return file_share.User{}, fmt.Errorf("error fetching user %d: %v", id, err)
	}

	return user, nil
}

func (us *UserService) GetFromSessionId(sessionId string) (file_share.User, error) {
	var user file_share.User

	if err := us.db.QueryRow("SELECT u.id, u.email, u.password, u.api_key FROM sessions s LEFT JOIN users u ON s.user_id = u.id WHERE session_id = ?", sessionId).Scan(&user.Id, &user.Email, &user.Password, &user.ApiKey); err != nil {
		if err == sql.ErrNoRows {
			return file_share.User{}, fmt.Errorf("session ID invalid")
		}

		return file_share.User{}, fmt.Errorf("error fetching user: %v", err)
	}

	return user, nil
}

func (us *UserService) GetFromCredentials(email, password string) (file_share.User, error) {
	user, err := us.GetFromEmail(email)

	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return file_share.User{}, fmt.Errorf("credentials invalid")
	}

	return user, nil
}

func (us *UserService) GetFromEmail(email string) (file_share.User, error) {
	var user file_share.User

	if err := us.db.QueryRow("SELECT id, email, password, api_key FROM users WHERE email = ?", email).Scan(&user.Id, &user.Email, &user.Password, &user.ApiKey); err != nil {
		if err == sql.ErrNoRows {
			return file_share.User{}, fmt.Errorf("credentials invalid")
		}

		return file_share.User{}, fmt.Errorf("error fetching user %s: %v", email, err)
	}

	return user, nil
}

func (us *UserService) Add(user file_share.User) (file_share.User, error) {
	result, err := us.db.Exec("INSERT INTO users (email, password, api_key) VALUES (?, ?, ?)", user.Email, user.Password, user.ApiKey)

	if err != nil {
		return file_share.User{}, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return file_share.User{}, err
	}

	return file_share.User{
		Id:    int(id),
		Email: user.Email,
	}, nil
}

func (us *UserService) UpdateEmail(user file_share.User, email string) error {
	if email == user.Email {
		return nil
	}

	_, err := us.db.Exec("UPDATE users SET email = ? WHERE id = ?", email, user.Id)

	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) UpdatePassword(user file_share.User, hash string) error {
	if hash == user.Password {
		return nil
	}

	_, err := us.db.Exec("UPDATE users SET password = ? WHERE id = ?", hash, user.Id)

	if err != nil {
		return err
	}

	return nil
}
