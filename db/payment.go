package db

import (
	"time"

	"gorm.io/gorm"
)

// For grocery payment
type Payment struct {
	gorm.Model
	AccountID  uint
	Account    *Account
	CategoryID uint
	Category   *Category
	Name       string
	Descr      string
	Note       string
	Items      []ItemBought
	Date       time.Time
}
