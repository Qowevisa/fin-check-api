package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Transfer struct {
	gorm.Model
	FromCardID              uint
	FromCard                *Card
	ToCardID                uint
	ToCard                  *Card
	Value                   uint64
	HaveDifferentCurrencies bool
	FromValue               uint64
	ToValue                 uint64
	Date                    time.Time
	UserID                  uint
	User                    *User
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
	ERROR_TRANSFER_FROMCARD_INSUFFICIENT_BALANCE = errors.New("FromCard's Balance is lower than Transfer's Value or FromValue")
	ERROR_TRANSFER_TOCARD_INVALID_USERID         = errors.New("Transfer's UserID and ToCard's UserID are not equal")
	ERROR_TRANSFER_CURRENCY_INCOSISTENCE         = errors.New("Transfer is using differenct currencies but FromValue or ToValue is not set approprietly")
)

func (t *Transfer) validateTransfer(fromCard, toCard *Card) error {
	if fromCard.UserID != t.UserID {
		return ERROR_TRANSFER_FROMCARD_INVALID_USERID
	}
	if toCard.UserID != t.UserID {
		return ERROR_TRANSFER_TOCARD_INVALID_USERID
	}
	diffCurs := fromCard.CurrencyID != toCard.CurrencyID
	if diffCurs && (t.FromValue == 0 || t.ToValue == 0) {
		return ERROR_TRANSFER_CURRENCY_INCOSISTENCE
	}
	if !diffCurs && fromCard.Balance < t.Value {
		return ERROR_TRANSFER_FROMCARD_INSUFFICIENT_BALANCE
	}
	if diffCurs && fromCard.Balance < t.FromValue {
		return ERROR_TRANSFER_FROMCARD_INSUFFICIENT_BALANCE
	}
	return nil
}

func (t *Transfer) BeforeCreate(tx *gorm.DB) error {
	fromCard := &Card{}
	if err := tx.Find(fromCard, t.FromCardID).Error; err != nil {
		return err
	}
	toCard := &Card{}
	if err := tx.Find(toCard, t.ToCardID).Error; err != nil {
		return err
	}
	if err := t.validateTransfer(fromCard, toCard); err != nil {
		return err
	}

	t.HaveDifferentCurrencies = fromCard.CurrencyID != toCard.CurrencyID
	// on same CurrencyID fromCard and toCard should use Value
	if !t.HaveDifferentCurrencies {
		fromCard.Balance -= t.Value
		if err := tx.Save(fromCard).Error; err != nil {
			return err
		}
		//
		toCard.Balance += t.Value
		if err := tx.Save(toCard).Error; err != nil {
			return err
		}
		return nil
	}
	// on DIFFERENT CurrencyID fromCard should use FromValue and toCard should use ToValue
	fromCard.Balance -= t.FromValue
	if err := tx.Save(fromCard).Error; err != nil {
		return err
	}
	//
	toCard.Balance += t.ToValue
	if err := tx.Save(toCard).Error; err != nil {
		return err
	}
	return nil
}

var (
	ERROR_TRANSFER_TOCARD_INSUFFICIENT_BALANCE = errors.New("Transfer's ToCard's Balance is lower than Value of ToValue of transfer")
)

func (t *Transfer) BeforeUpdate(tx *gorm.DB) error {
	var original *Transfer
	if err := tx.Find(original, t.ID).Error; err != nil {
		return err
	}
	var origFromCard *Card
	if err := tx.Find(origFromCard, original.FromCardID).Error; err != nil {
		return err
	}
	var origToCard *Card
	if err := tx.Find(origToCard, original.ToCardID).Error; err != nil {
		return err
	}

	diffCurs := origFromCard.CurrencyID != origToCard.CurrencyID
	if !diffCurs {
		origFromCard.Balance += original.Value
		if origToCard.Balance < original.Value {
			return ERROR_TRANSFER_TOCARD_INSUFFICIENT_BALANCE
		}
		origToCard.Balance -= original.Value
	} else {
		if original.FromValue == 0 || original.ToValue == 0 {
			return ERROR_TRANSFER_CURRENCY_INCOSISTENCE
		}
		origFromCard.Balance += original.FromValue
		if origToCard.Balance < original.ToValue {
			return ERROR_TRANSFER_TOCARD_INSUFFICIENT_BALANCE
		}
		origToCard.Balance -= original.ToValue
	}
	if err := tx.Save(origFromCard).Error; err != nil {
		return err
	}
	if err := tx.Save(origToCard).Error; err != nil {
		return err
	}
	//
	fromCard := &Card{}
	if err := tx.Find(fromCard, t.FromCardID).Error; err != nil {
		return err
	}
	toCard := &Card{}
	if err := tx.Find(toCard, t.ToCardID).Error; err != nil {
		return err
	}
	if err := t.validateTransfer(fromCard, toCard); err != nil {
		return err
	}

	t.HaveDifferentCurrencies = fromCard.CurrencyID != toCard.CurrencyID
	// on same CurrencyID fromCard and toCard should use Value
	if !t.HaveDifferentCurrencies {
		fromCard.Balance -= t.Value
		if err := tx.Save(fromCard).Error; err != nil {
			return err
		}
		//
		toCard.Balance += t.Value
		if err := tx.Save(toCard).Error; err != nil {
			return err
		}
		return nil
	}
	// on DIFFERENT CurrencyID fromCard should use FromValue and toCard should use ToValue
	fromCard.Balance -= t.FromValue
	if err := tx.Save(fromCard).Error; err != nil {
		return err
	}
	//
	toCard.Balance += t.ToValue
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
	toCard := &Card{}
	if err := tx.Find(toCard, t.ToCardID).Error; err != nil {
		return err
	}
	diffCurs := fromCard.CurrencyID != toCard.CurrencyID
	if !diffCurs {
		fromCard.Balance += t.Value
		if toCard.Balance < t.Value {
			return ERROR_TRANSFER_TOCARD_INSUFFICIENT_BALANCE
		}
		toCard.Balance -= t.Value
		return nil
	} else {
		if t.FromValue == 0 || t.ToValue == 0 {
			return ERROR_TRANSFER_CURRENCY_INCOSISTENCE
		}
		fromCard.Balance += t.FromValue
		if toCard.Balance < t.ToValue {
			return ERROR_TRANSFER_TOCARD_INSUFFICIENT_BALANCE
		}
		toCard.Balance -= t.ToValue
	}
	if err := tx.Save(fromCard).Error; err != nil {
		return err
	}
	if err := tx.Save(toCard).Error; err != nil {
		return err
	}
	return nil
}
