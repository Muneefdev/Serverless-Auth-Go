package database

import "lambda-func/types"

type UserStore interface {
	DoesUserExists(username string) (bool, error)
	InsertUser(user *types.User) error
	GetUserByName(username string) (*types.User, error)
}
