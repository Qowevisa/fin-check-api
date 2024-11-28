package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	gorm.Model
	CardID  uint
	Card    *Card
	Value   uint64
	Comment string
	Date    time.Time
	UserID  uint
	User    *User
	TypeID  uint
	Type    *Type
}

type Helper_ExpenseBulk struct {
	PropagateCardID  bool
	CardID           uint
	PropagateTypeID  bool
	TypeID           uint
	PropagateValue   bool
	Value            uint64
	PropagateComment bool
	Comment          string
	PropagateDate    bool
	Date             time.Time
	UserID           uint
}

// {{{ Helper_ExpenseBulk.CreateExpenseFromChild I'm not proud of this code
func (he *Helper_ExpenseBulk) CreateExpenseFromChild(c Expense) *Expense {
	var cardID uint
	var typeID uint
	var value uint64
	var comment string
	var date time.Time
	if he.PropagateCardID {
		cardID = he.CardID
	} else {
		cardID = c.CardID
	}
	if he.PropagateTypeID {
		typeID = he.TypeID
	} else {
		typeID = c.TypeID
	}
	if he.PropagateValue {
		value = he.Value
	} else {
		value = c.Value
	}
	if he.PropagateComment {
		comment = he.Comment
	} else {
		comment = c.Comment
	}
	if he.PropagateValue {
		value = he.Value
	} else {
		value = c.Value
	}
	if he.PropagateDate {
		date = he.Date
	} else {
		date = c.Date
	}
	return &Expense{
		CardID:  cardID,
		TypeID:  typeID,
		Value:   value,
		Comment: comment,
		Date:    date,
		UserID:  he.UserID,
	}
}

// }}}

// Implements db.UserIdentifiable:1
func (e Expense) GetID() uint {
	return e.ID
}

// Implements db.UserIdentifiable:2
func (e Expense) GetUserID() uint {
	return e.UserID
}

// Implements db.UserIdentifiable:3
func (e *Expense) SetUserID(id uint) {
	e.UserID = id
}

var (
	ERROR_EXPENSE_INVALID_CARD_USERID       = errors.New("Expense's `UserID` and Card's `UserID` are not equal")
	ERROR_EXPENSE_CARD_INSUFFICIENT_BALANCE = errors.New("Card's `Balance` is lower than Expense's Value")
	ERROR_EXPENSE_INVALID_TYPE_USERID       = errors.New("Expense's `UserID` and Type's `UserID` are not equal")
)

func (e *Expense) BeforeCreate(tx *gorm.DB) error {
	card := &Card{}
	if err := tx.Find(card, e.CardID).Error; err != nil {
		return err
	}
	if card.UserID != e.UserID {
		return ERROR_EXPENSE_INVALID_CARD_USERID
	}
	if card.Balance < e.Value {
		return ERROR_EXPENSE_CARD_INSUFFICIENT_BALANCE
	}
	card.Balance -= e.Value
	if err := tx.Save(card).Error; err != nil {
		return err
	}

	typ := &Type{}
	if err := tx.Find(typ, e.TypeID).Error; err != nil {
		return err
	}
	if typ.UserID != e.UserID {
		return ERROR_EXPENSE_INVALID_TYPE_USERID
	}
	return nil
}

func (e *Expense) BeforeUpdate(tx *gorm.DB) (err error) {
	var original Expense
	if err := tx.Model(&Expense{}).Select("card_id", "value").Where("id = ?", e.ID).First(&original).Error; err != nil {
		return err
	}
	if original.CardID != 0 {
		oldCard := &Card{}
		if err := tx.Find(oldCard, original.CardID).Error; err != nil {
			return err
		}
		if oldCard.UserID != e.UserID {
			return ERROR_EXPENSE_INVALID_CARD_USERID
		}
		oldCard.Balance += original.Value
		if err := tx.Save(oldCard).Error; err != nil {
			return err
		}
	}

	if e.CardID != 0 {
		newCard := &Card{}
		if err := tx.Find(newCard, e.CardID).Error; err != nil {
			return err
		}
		if newCard.UserID != e.UserID {
			return ERROR_EXPENSE_INVALID_CARD_USERID
		}
		if newCard.Balance < e.Value {
			return ERROR_EXPENSE_CARD_INSUFFICIENT_BALANCE
		}
		newCard.Balance -= e.Value
		if err := tx.Save(newCard).Error; err != nil {
			return err
		}
	}
	typ := &Type{}
	if err := tx.Find(typ, e.TypeID).Error; err != nil {
		return err
	}
	if typ.UserID != e.UserID {
		return ERROR_EXPENSE_INVALID_TYPE_USERID
	}
	return nil
}

func (e *Expense) AfterDelete(tx *gorm.DB) (err error) {
	card := &Card{}
	if err := tx.Find(card, e.CardID).Error; err != nil {
		return err
	}
	if card.UserID != e.UserID {
		return ERROR_EXPENSE_INVALID_CARD_USERID
	}
	card.Balance += e.Value
	if err := tx.Save(card).Error; err != nil {
		return err
	}
	return nil
}
