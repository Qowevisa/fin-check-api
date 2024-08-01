package db

import "gorm.io/gorm"

type ItemBought struct {
	gorm.Model
	ItemID    uint
	Item      *Item
	Quantity  uint
	PaymentID uint
	Payment   *Payment
}
