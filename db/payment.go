package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// For grocery payment
type Payment struct {
	gorm.Model
	CardID     uint
	Card       *Card
	CategoryID uint
	Category   *Category
	UserID     uint
	User       *User
	Title      string
	Descr      string
	Note       string
	Items      []ItemBought `gorm:"constraint:OnDelete:CASCADE;"`
	Date       time.Time
}

func (p Payment) __internalBelogingToPayment() {}

var (
	ERROR_PAYMENT_INVALID_CARD_USERID     = errors.New("Payment's `UserID` and Card's `UserID` are not equal")
	ERROR_PAYMENT_INVALID_CATEGORY_USERID = errors.New("Payment's `UserID` and Category's `UserID` are not equal")
)

func (p *Payment) BeforeSave(tx *gorm.DB) error {
	paymentCard := &Card{}
	if err := tx.Find(paymentCard, p.CardID).Error; err != nil {
		return err
	}
	if paymentCard.UserID != p.UserID {
		return ERROR_PAYMENT_INVALID_CARD_USERID
	}
	paymentCategory := &Category{}
	if err := tx.Find(paymentCategory, p.CategoryID).Error; err != nil {
		return err
	}
	if paymentCategory.UserID != p.UserID {
		return ERROR_PAYMENT_INVALID_CATEGORY_USERID
	}
	return nil
}
