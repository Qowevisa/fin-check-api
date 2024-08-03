package db

import (
	"errors"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name string
	// Parent is used as a infinite sub-category structure
	ParentID uint
	Parent   *Category
	UserID   uint
	User     *User
}

var (
	ERROR_CATEGORY_PARENT_NOT_FOUND  = errors.New("ParentID is invalid for user")
	ERROR_CATEGORY_NAME_NOT_UNIQUE   = errors.New("Name for Category have to be unique for user")
	ERROR_CATEGORY_USER_ID_NOT_EQUAL = errors.New("ParentID is invalid for user")
)

func (c *Category) BeforeSave(tx *gorm.DB) error {
	if c.ParentID != 0 {
		var parent Category
		if err := tx.Find(&parent, c.ID).Error; err != nil {
			return err
		}
		if parent.ID == 0 {
			return ERROR_CATEGORY_PARENT_NOT_FOUND
		}
		if parent.UserID != c.UserID {
			return ERROR_CATEGORY_USER_ID_NOT_EQUAL
		}
	}
	var dup Category
	if err := tx.Find(&dup, Category{Name: c.Name, UserID: c.UserID}).Error; err != nil {
		return err
	}
	if dup.ID != 0 {
		return ERROR_CATEGORY_NAME_NOT_UNIQUE
	}
	return nil
}
