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

var (
	ERROR_CARD_NAME_EMPTY      = errors.New("Card's name can't be empty")
	ERROR_CARD_NAME_NOT_UNIQUE = errors.New("Card's name have to be unique for user")
)

func (c *Card) BeforeSave(tx *gorm.DB) error {
	if c.Name == "" {
		return ERROR_CARD_NAME_EMPTY
	}

	var dup Card
	if err := tx.Find(&dup, Card{Name: c.Name}).Error; err != nil {
		return err
	}

	if dup.ID != 0 {
		return ERROR_CARD_NAME_NOT_UNIQUE
	}

	return nil
}
