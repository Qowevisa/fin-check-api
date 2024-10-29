package db

import "gorm.io/gorm"

type Type struct {
	gorm.Model
	Name    string
	Comment string
	Color   string
}
