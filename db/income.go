package db

import (
	"time"

	"gorm.io/gorm"
)

type Income struct {
	gorm.Model
	CardID uint
	Card   *Card
	Value  uint64
	Date   time.Time
}
