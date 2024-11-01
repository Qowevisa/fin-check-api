package db

import (
	"time"

	"gorm.io/gorm"
)

type Debt struct {
	gorm.Model
	CardID   uint
	Card     *Card
	Comment  string
	Value    uint64
	IOwe     bool
	Date     time.Time
	DateEnd  time.Time
	Finished bool
	UserID   uint
	User     *User
}

// Implements db.UserIdentifiable:1
func (d Debt) GetID() uint {
	return d.ID
}

// Implements db.UserIdentifiable:2
func (d Debt) GetUserID() uint {
	return d.UserID
}

// Implements db.UserIdentifiable:3
func (d *Debt) SetUserID(id uint) {
	d.UserID = id
}
