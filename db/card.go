package db

import "gorm.io/gorm"

// Card can be either card or wallet
type Card struct {
	gorm.Model
	Name           string
	Value          uint64
	HaveCreditLine bool
	CreditLine     uint64
}
