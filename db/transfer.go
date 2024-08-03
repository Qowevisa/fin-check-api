package db

import "gorm.io/gorm"

type Transfer struct {
	gorm.Model
	FromCardID uint
	FromCard   *Card
	ToCardID   uint
	ToCard     *Card
	Value      uint64
}
