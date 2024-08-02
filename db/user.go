package db

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

var (
	ERROR_USER_USERNAME_NOT_UNIQUE = errors.New("Username already persists in database")
)

func (u *User) BeforeCreate(ctx *gorm.DB) error {
	var duplicate User
	err := ctx.Find(&duplicate, User{Username: u.Username}).Error
	if err != nil {
		return err
	}
	if duplicate.ID != 0 {
		return ERROR_USER_USERNAME_NOT_UNIQUE
	}
	return nil
}
