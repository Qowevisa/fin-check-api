package db

import "gorm.io/gorm"

type ItemBought struct {
	gorm.Model
	ItemID      uint
	Item        *Item
	PaymentID   uint
	Payment     *Payment
	TypeID      uint
	Type        *Type
	Quantity    uint
	TotalCost   uint64
	MetricType  uint8
	MetricValue uint64
}

func (i ItemBought) __internalBelogingToPayment() {}
