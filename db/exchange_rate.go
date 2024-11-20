package db

import "gorm.io/gorm"

type ExchangeRate struct {
	gorm.Model
	FromCurrID uint
	FromCurr   *Currency
	From       uint64
	ToCurrID   uint
	ToCurr     *Currency
	To         uint64
	Rate       uint64
}
