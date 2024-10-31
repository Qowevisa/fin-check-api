package db

import (
	"errors"

	"gorm.io/gorm"
)

// Card can be either card or wallet
type Card struct {
	gorm.Model
	Name           string
	Value          uint64
	HaveCreditLine bool
	CreditLine     uint64
	UserID         uint
	User           *User
}

// Implements db.UserIdentifiable:1
func (c Card) GetID() uint {
	return c.ID
}

// Implements db.UserIdentifiable:2
func (c Card) GetUserID() uint {
	return c.UserID
}

// Implements db.UserIdentifiable:3
func (c *Card) SetUserID(id uint) {
	c.UserID = id
}

var (
	ERROR_CARD_NAME_EMPTY      = errors.New("The 'Name' field for 'Card' cannot be empty")
	ERROR_CARD_NAME_NOT_UNIQUE = errors.New("The 'Name' field for 'Card' have to be unique for user")
)

func (c *Card) BeforeSave(tx *gorm.DB) error {
	if c.Name == "" {
		return ERROR_CARD_NAME_EMPTY
	}

	var dup Card
	if err := tx.Find(&dup, Card{Name: c.Name, UserID: c.UserID}).Error; err != nil {
		return err
	}

	if c.ID != dup.ID && dup.ID != 0 {
		return ERROR_CARD_NAME_NOT_UNIQUE
	}

	return nil
}
