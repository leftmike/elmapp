package model

import (
	"errors"
	"fmt"
)

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	password string
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

var (
	userByEmail    = map[string]*User{}
	userByUsername = map[string]*User{}
)

func LoginUser(email, password string) *User {
	user, ok := userByEmail[email]
	if !ok || user.password != password {
		fmt.Printf("login user: %s %s FAILED\n", email, password)
		return nil
	}
	token, err := newToken(user)
	if err != nil {
		fmt.Printf("login user: new token: %s\n", err)
		return nil
	}
	user.Token = token
	fmt.Printf("login user: %s %s SUCCEEDED\n", email, password)
	return user
}

func RegisterUser(username, email, password string) (*User, error) {
	if _, ok := userByUsername[username]; ok {
		fmt.Printf("register user: username %s already registered\n", username)
		return nil, errors.New("username already registered")
	}
	if _, ok := userByEmail[email]; ok {
		fmt.Printf("register user: email address %s already registered\n", email)
		return nil, errors.New("email address already registered")
	}
	fmt.Printf("register user: %s %s %s SUCCEEDED\n", username, email, password)
	user := &User{
		Username: username,
		Email:    email,
		password: password,
	}
	userByUsername[username] = user
	userByEmail[email] = user
	return user, nil
}
