package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Income struct {
	gorm.Model
	CardID  uint
	Card    *Card
	Value   uint64
	Comment string
	Date    time.Time
	UserID  uint
	User    *User
}

// Implements db.UserIdentifiable:1
func (e Income) GetID() uint {
	return e.ID
}

// Implements db.UserIdentifiable:2
func (e Income) GetUserID() uint {
	return e.UserID
}

// Implements db.UserIdentifiable:3
func (e *Income) SetUserID(id uint) {
	e.UserID = id
}

var (
	ERROR_INCOME_INVALID_USERID            = errors.New("Income's `UserID` and Card's `UserID` are not equal")
	ERROR_INCOME_CARD_INSUFFICIENT_BALANCE = errors.New("Card's `Balance` is lower than Income's Value")
	ERROR_INCOME_INVALID_TYPE_USERID       = errors.New("Income's `UserID` and Type's `UserID` are not equal")
)

func (i *Income) BeforeCreate(tx *gorm.DB) error {
	card := &Card{}
	if err := tx.Find(card, i.CardID).Error; err != nil {
		return err
	}
	if card.UserID != i.UserID {
		return ERROR_INCOME_INVALID_USERID
	}
	card.Balance += i.Value
	if err := tx.Save(card).Error; err != nil {
		return err
	}

	return nil
}

func (i *Income) BeforeUpdate(tx *gorm.DB) (err error) {
	var original Income
	if err := tx.Model(&Income{}).Select("card_id", "value").Where("id = ?", i.ID).First(&original).Error; err != nil {
		return err
	}
	if original.CardID != 0 {
		oldCard := &Card{}
		if err := tx.Find(oldCard, original.CardID).Error; err != nil {
			return err
		}
		if oldCard.UserID != i.UserID {
			return ERROR_INCOME_INVALID_USERID
		}
		if oldCard.Balance < original.Value {
			return ERROR_INCOME_CARD_INSUFFICIENT_BALANCE
		}
		oldCard.Balance -= original.Value
		if err := tx.Save(oldCard).Error; err != nil {
			return err
		}
	}

	if i.CardID != 0 {
		newCard := &Card{}
		if err := tx.Find(newCard, i.CardID).Error; err != nil {
			return err
		}
		if newCard.UserID != i.UserID {
			return ERROR_INCOME_INVALID_USERID
		}
		newCard.Balance += i.Value
		if err := tx.Save(newCard).Error; err != nil {
			return err
		}
	}
	return nil
}

func (e *Income) AfterDelete(tx *gorm.DB) (err error) {
	card := &Card{}
	if err := tx.Find(card, e.CardID).Error; err != nil {
		return err
	}
	if card.UserID != e.UserID {
		return ERROR_INCOME_INVALID_USERID
	}
	card.Balance -= e.Value
	if err := tx.Save(card).Error; err != nil {
		return err
	}
	return nil
}
