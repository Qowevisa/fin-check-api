package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Transfer struct {
	gorm.Model
	FromCardID uint
	FromCard   *Card
	ToCardID   uint
	ToCard     *Card
	Value      uint64
	Date       time.Time
	UserID     uint
	User       *User
}

// Implements db.UserIdentifiable:1
func (t Transfer) GetID() uint {
	return t.ID
}

// Implements db.UserIdentifiable:2
func (t Transfer) GetUserID() uint {
	return t.UserID
}

// Implements db.UserIdentifiable:3
func (t *Transfer) SetUserID(id uint) {
	t.UserID = id
}

var (
	ERROR_TRANSFER_FROMCARD_INVALID_USERID       = errors.New("Transfer's UserID and FromCard's UserID are not equal")
	ERROR_TRANSFER_FROMCARD_INSUFFICIENT_BALANCE = errors.New("FromCard's Balance is lower than Transfer's Value")
	ERROR_TRANSFER_TOCARD_INVALID_USERID         = errors.New("Transfer's UserID and ToCard's UserID are not equal")
)

func (t *Transfer) BeforeCreate(tx *gorm.DB) error {
	fromCard := &Card{}
	if err := tx.Find(fromCard, t.FromCardID).Error; err != nil {
		return err
	}
	if fromCard.UserID != t.UserID {
		return ERROR_TRANSFER_FROMCARD_INVALID_USERID
	}
	if fromCard.Balance < t.Value {
		return ERROR_TRANSFER_FROMCARD_INSUFFICIENT_BALANCE
	}
	fromCard.Balance -= t.Value
	if err := tx.Save(fromCard).Error; err != nil {
		return err
	}
	//
	toCard := &Card{}
	if err := tx.Find(toCard, t.ToCardID).Error; err != nil {
		return err
	}
	if toCard.UserID != t.UserID {
		return ERROR_TRANSFER_TOCARD_INVALID_USERID
	}
	toCard.Balance += t.Value
	if err := tx.Save(toCard).Error; err != nil {
		return err
	}
	return nil
}

var (
	ERROR_TRANSFER_FROMCARD_IDZERO             = errors.New("Transfer's FromCardID is zero")
	ERROR_TRANSFER_TOCARD_IDZERO               = errors.New("Transfer's ToCardID is zero")
	ERROR_TRANSFER_TOCARD_INSUFFICIENT_BALANCE = errors.New("Transfer's ToCard's Balance is lower than value of transfer")
)

func (t *Transfer) BeforeUpdate(tx *gorm.DB) error {
	var original Transfer
	if err := tx.Find(&original, t.ID).Error; err != nil {
		return err
	}
	if original.FromCardID == 0 {
		return ERROR_TRANSFER_FROMCARD_IDZERO
	}
	var origFromCard Card
	if err := tx.Find(&origFromCard, original.FromCardID).Error; err != nil {
		return err
	}
	origFromCard.Balance += original.Value
	if err := tx.Save(origFromCard).Error; err != nil {
		return err
	}

	if original.ToCardID == 0 {
		return ERROR_TRANSFER_FROMCARD_IDZERO
	}
	var origToCard Card
	if err := tx.Find(&origToCard, original.ToCardID).Error; err != nil {
		return err
	}
	if origToCard.Balance < original.Value {
		return ERROR_TRANSFER_TOCARD_INSUFFICIENT_BALANCE
	}
	origToCard.Balance -= original.Value
	if err := tx.Save(origToCard).Error; err != nil {
		return err
	}
	//
	fromCard := &Card{}
	if err := tx.Find(fromCard, t.FromCardID).Error; err != nil {
		return err
	}
	if fromCard.UserID != t.UserID {
		return ERROR_TRANSFER_FROMCARD_INVALID_USERID
	}
	if fromCard.Balance < t.Value {
		return ERROR_TRANSFER_FROMCARD_INSUFFICIENT_BALANCE
	}
	fromCard.Balance -= t.Value
	if err := tx.Save(fromCard).Error; err != nil {
		return err
	}
	//
	toCard := &Card{}
	if err := tx.Find(toCard, t.ToCardID).Error; err != nil {
		return err
	}
	if toCard.UserID != t.UserID {
		return ERROR_TRANSFER_TOCARD_INVALID_USERID
	}
	toCard.Balance += t.Value
	if err := tx.Save(toCard).Error; err != nil {
		return err
	}
	return nil
}

func (t *Transfer) AfterDelete(tx *gorm.DB) error {
	fromCard := &Card{}
	if err := tx.Find(fromCard, t.FromCardID).Error; err != nil {
		return err
	}
	fromCard.Balance += t.Value
	if err := tx.Save(fromCard).Error; err != nil {
		return err
	}
	toCard := &Card{}
	if err := tx.Find(toCard, t.ToCardID).Error; err != nil {
		return err
	}
	if toCard.Balance < t.Value {
		return ERROR_TRANSFER_TOCARD_INSUFFICIENT_BALANCE
	}
	toCard.Balance -= t.Value
	if err := tx.Save(toCard).Error; err != nil {
		return err
	}
	return nil
}
