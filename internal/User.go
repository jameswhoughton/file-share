package user

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	password string
	Email    string
}

type Form struct {
	Email    string
	Password string
	ApiKey   string
}

type Model struct {
	db *sql.DB
}

func NewUserModel(db *sql.DB) Model {
	return Model{db}
}

func (um *Model) Get(id int) (User, error) {
	var user User

	if err := um.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user); err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user %d not found", id)
		}

		return User{}, fmt.Errorf("error fetching user %d: %v", id, err)
	}

	return user, nil
}

func (um *Model) GetWithCredentials(email, password string) (User, error) {
	var user User

	if err := um.db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user.Id, &user.Email, &user.password); err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("credentials invalid")
		}

		return User{}, fmt.Errorf("error fetching user %s: %v", email, err)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.password), []byte(password))

	if err != nil {
		return User{}, fmt.Errorf("credentials invalid")
	}

	return user, nil
}

func (um *Model) Add(user Form) (User, error) {
	result, err := um.db.Exec("INSERT INTO users (email, password, api_key) VALUES (?, ?, ?)", user.Email, user.Password, user.ApiKey)

	if err != nil {
		return User{}, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return User{}, err
	}

	return User{
		Id:    int(id),
		Email: user.Email,
	}, nil
}
