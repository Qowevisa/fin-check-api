package db

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var udb *gorm.DB
var conMu sync.Mutex

var (
	ERROR_DB_NOT_INIT = errors.New("Database connection is not initialized")
)

func Connect() *gorm.DB {
	conMu.Lock()
	defer conMu.Unlock()
	if udb != nil {
		return udb
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)
	gormDB, err := gorm.Open(sqlite.Open("gonuts.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Panic(err)
	}
	newUDB := gormDB
	gormDB.AutoMigrate(&Card{})
	gormDB.AutoMigrate(&Category{})
	gormDB.AutoMigrate(&Item{})
	gormDB.AutoMigrate(&ItemPrice{})
	gormDB.AutoMigrate(&Payment{})
	gormDB.AutoMigrate(&ItemBought{})
	gormDB.AutoMigrate(&Income{})
	gormDB.AutoMigrate(&Debt{})
	gormDB.AutoMigrate(&Transfer{})
	gormDB.AutoMigrate(&User{})
	return newUDB
}
