package db

import (
	"time"

	"gorm.io/gorm"
)

type Debt struct {
	gorm.Model
	CardID   uint
	Card     *Card
	Value    uint64
	IOwe     bool
	Date     time.Time
	DateEnd  time.Time
	Finished bool
}
