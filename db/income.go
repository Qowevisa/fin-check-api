package db

import (
	"time"

	"gorm.io/gorm"
)

type Income struct {
	gorm.Model
	AccountID uint
	Account   *Account
	Value     uint64
	Date      time.Time
}
