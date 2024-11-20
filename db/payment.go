package db

import (
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
