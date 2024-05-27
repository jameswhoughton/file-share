package main

import (
	"database/sql"
	"fmt"
)

type User struct {
	id    int
	email string
}

type UserForm struct {
	email    string
	password string
	apiKey   string
}

type UserModel struct {
	db *sql.DB
}

func (um *UserModel) Get(id int) (User, error) {
	var user User

	if err := um.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user); err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user %d not found", id)
		}

		return User{}, fmt.Errorf("error fetching user %d: %v", id, err)
	}

	return user, nil
}

func (um *UserModel) GetWithCredentials(email, password string) (User, error) {
	var user User

	if err := um.db.QueryRow("SELECT * FROM users WHERE email = ? AND password = ?", email, password).Scan(&user); err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("credentials invalid")
		}

		return User{}, fmt.Errorf("error fetching user %s: %v", email, err)
	}

	return user, nil
}

func (um *UserModel) Add(user UserForm) (User, error) {
	result, err := um.db.Exec("INSERT INTO users (email, password, api_key) VALUES (?, ?, ?)", user.email, user.password, user.apiKey)

	if err != nil {
		return User{}, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return User{}, err
	}

	return User{
		id:    int(id),
		email: user.email,
	}, nil
}
