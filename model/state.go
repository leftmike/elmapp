package model

import (
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

func RegisterUser(username, email, password string) (*User, []string) {
	var ret []string

	if _, ok := userByUsername[username]; ok {
		fmt.Printf("register user: username %s already registered\n", username)
		ret = append(ret, "username has already been taken")
	}
	if _, ok := userByEmail[email]; ok {
		fmt.Printf("register user: email address %s already registered\n", email)
		ret = append(ret, "email has already been taken")
	}
	if len(password) < 8 {
		ret = append(ret, "password is too short (minimum is 8 characters)")
	}
	if ret != nil {
		return nil, ret
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

var (
	defaultUsers = []User{
		{Username: "setup", Email: "setup@setup.com", password: "default"},
		{Username: "mike", Email: "mike@mike.com", password: "password"},
		{Username: "test", Email: "test@test.com", password: "test"},
	}
)

func init() {
	for idx, u := range defaultUsers {
		userByUsername[u.Username] = &defaultUsers[idx]
		userByEmail[u.Email] = &defaultUsers[idx]
	}
}
