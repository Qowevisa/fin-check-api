package db

import "gorm.io/gorm"

type Metric struct {
	gorm.Model
	Value uint8
	Name  string
	Short string
}

func (m Metric) GetID() uint {
	return m.ID
}
