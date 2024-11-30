package db

import (
	"errors"

	"gorm.io/gorm"
)

type SettingsTypeFilter struct {
	gorm.Model
	UserID uint
	User   *User
	TypeID uint
	Type   *Type
}

// so `SETTINGS_TYPE` is the NAME of struct that has such error
// second `TYPE` means SettingsTypeFilter's Type fields are causing errors
var (
	ERROR_SETTINGS_TYPE_TYPEID_ZERO    = errors.New("SettingsTypeFilter's `TypeID` is zero")
	ERROR_SETTINGS_TYPE_USERID_INVALID = errors.New("SettingsTypeFilter's `UserID` and Type's `UserID` are not equal")
)

func (st *SettingsTypeFilter) BeforeSave(tx *gorm.DB) error {
	if st.TypeID == 0 {
		return ERROR_SETTINGS_TYPE_TYPEID_ZERO
	}
	dbType := &Type{}
	if err := tx.Find(dbType, st.TypeID).Error; err != nil {
		return err
	}
	if dbType.UserID != st.UserID {
		return ERROR_SETTINGS_TYPE_USERID_INVALID
	}
	return nil
}
