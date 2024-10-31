package db

import (
	"errors"

	"gorm.io/gorm"
)

type Type struct {
	gorm.Model
	Name    string
	Comment string
	Color   string
	UserID  uint
	User    *User
}

// Implements db.UserIdentifiable:1
func (t Type) GetID() uint {
	return t.ID
}

// Implements db.UserIdentifiable:2
func (t Type) GetUserID() uint {
	return t.UserID
}

// Implements db.UserIdentifiable:3
func (t *Type) SetUserID(id uint) {
	t.UserID = id
}

var (
	ERROR_TYPE_NAME_EMPTY      = errors.New("The 'Name' field of 'Type' cannot be empty")
	ERROR_TYPE_NAME_NOT_UNIQUE = errors.New("The 'Name' field of 'Type' have to be unique for user")
)

func (t *Type) BeforeSave(tx *gorm.DB) error {
	if t.Name == "" {
		return ERROR_TYPE_NAME_EMPTY
	}

	var dup Type
	if err := tx.Find(&dup, Type{Name: t.Name, UserID: t.UserID}).Error; err != nil {
		return err
	}

	if t.ID != dup.ID && dup.ID != 0 {
		return ERROR_TYPE_NAME_NOT_UNIQUE
	}

	return nil
}
