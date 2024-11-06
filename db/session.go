package db

import "time"

type Session struct {
	ID       string `gorm:"primaryKey"`
	UserID   uint
	User     *User
	ExpireAt time.Time
}
