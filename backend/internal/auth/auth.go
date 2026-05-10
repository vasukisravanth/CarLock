package auth

import (
    "errors"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    Username string
    Password string
}

var users = map[string]string{}

func Register(username, password string) error {
    if _, exists := users[username]; exists {
        return errors.New("username already exists")
    }
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    users[username] = string(hashedPassword)
    return nil
}

func Authenticate(username, password string) (bool, error) {
    hashedPassword, exists := users[username]
    if !exists {
        return false, nil
    }
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        return false, err
    }
    return true, nil
}