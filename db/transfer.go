package db

import (
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
