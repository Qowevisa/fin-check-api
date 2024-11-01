package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type ItemPrice struct {
	gorm.Model
	ItemID    uint
	Item      *Item
	Price     uint64
	ValidFrom time.Time
	ValidTo   time.Time
}

var (
	ERROR_ITEMPRICE_VALID_FROM_ERR = errors.New("ValidFrom shall be initiated when created")
)

func (i *ItemPrice) BeforeCreate(tx *gorm.DB) error {
	if i.ValidFrom.IsZero() {
		i.ValidFrom = time.Now()
	}
	return nil
}
