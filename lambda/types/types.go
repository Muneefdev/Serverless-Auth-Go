package types

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser(registerUser User) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 12)
	if err != nil {
		return nil, err
	}

	return &User{
		Username: registerUser.Username,
		Password: string(passwordHash),
	}, nil
}

func ValidatePassword(passwordHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}
