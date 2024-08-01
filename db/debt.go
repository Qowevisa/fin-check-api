package db

import (
	"time"

	"gorm.io/gorm"
)

type Debt struct {
	gorm.Model
	AccountID uint
	Account   *Account
	Value     uint64
	IOwe      bool
	Date      time.Time
	DateEnd   time.Time
	Finished  bool
}
