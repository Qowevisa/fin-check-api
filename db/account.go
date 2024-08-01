package db

import "gorm.io/gorm"

// Account can be either card or wallet
type Account struct {
	gorm.Model
	Name           string
	Value          uint64
	HaveCreditLine bool
	CreditLine     uint64
}
