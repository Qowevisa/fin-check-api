package db

import "gorm.io/gorm"

type Currency struct {
	gorm.Model
	Name    string
	ISOName string
	Symbol  string
}

func (c Currency) GetID() uint {
	return c.ID
}
