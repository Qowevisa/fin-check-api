package db

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string
	// Parent is used as a infinite sub-category structure
	ParentID uint
	Parent   *Category
	UserID   uint
	User     *User
}
