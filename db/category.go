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

// Implements db.UserIdentifiable:1
func (c Category) GetID() uint {
	return c.ID
}

// Implements db.UserIdentifiable:2
func (c Category) GetUserID() uint {
	return c.UserID
}

// Implements db.UserIdentifiable:3
func (c *Category) SetUserID(id uint) {
	c.UserID = id
}

var (
	ERROR_CATEGORY_PARENT_NOT_FOUND  = errors.New("Can't find Category with ParentID for user")
	ERROR_CATEGORY_NAME_NOT_UNIQUE   = errors.New("Name for Category have to be unique for user")
	ERROR_CATEGORY_USER_ID_NOT_EQUAL = errors.New("ParentID is invalid for user")
	ERROR_CATEGORY_SELF_REFERENCING  = errors.New("Category can't set itself as a parent")
)

func (c *Category) BeforeSave(tx *gorm.DB) error {
	if c.ParentID != 0 {
		var parent Category
		if err := tx.Find(&parent, c.ParentID).Error; err != nil {
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
	if c.ID != dup.ID && dup.ID != 0 {
		return ERROR_CATEGORY_NAME_NOT_UNIQUE
	}
	return nil
}

func (c *Category) AfterCreate(tx *gorm.DB) error {
	if c.ParentID == c.ID {
		return ERROR_CATEGORY_SELF_REFERENCING
	}
	return nil
}
