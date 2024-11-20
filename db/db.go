package db

import (
	"errors"
	"fmt"
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

func Init() error {
	dbc := Connect()
	// Seeds
	if err := initStateOfDb(dbc); err != nil {
		return fmt.Errorf("initStateOfDb: %w", err)
	}
	return nil
}

func Connect() *gorm.DB {
	conMu.Lock()
	defer conMu.Unlock()
	if udb != nil {
		return udb
	}
	logFile, err := os.OpenFile("db.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	newLogger := logger.New(
		log.New(logFile, "\r\n", log.LstdFlags),
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
	gormDB.AutoMigrate(&Type{})
	gormDB.AutoMigrate(&Session{})
	gormDB.AutoMigrate(&Expense{})
	gormDB.AutoMigrate(&Metric{})
	gormDB.AutoMigrate(&Currency{})
	return newUDB
}

var (
	CANT_FIND_METRIC = errors.New("Can't find proper metrics in database")
)

func initMetrics(tx *gorm.DB) error {
	var metrics []Metric
	if err := tx.Find(&metrics).Error; err != nil {
		return err
	}
	metricsThatNeeded := []*Metric{
		&Metric{Name: "None", Short: "pcs", Value: 0},
		&Metric{Name: "Gram", Short: "g", Value: 1},
		&Metric{Name: "Kilogram", Short: "kg", Value: 2},
		&Metric{Name: "Liter", Short: "l", Value: 3},
	}
	if len(metrics) == 0 {
		for _, m := range metricsThatNeeded {
			if err := tx.Create(m).Error; err != nil {
				return err
			}
		}
		return nil
	}
	for _, m := range metricsThatNeeded {
		tmp := &Metric{}
		if err := tx.Find(tmp, m).Error; err != nil {
			return err
		}
		if tmp.ID == 0 {
			return CANT_FIND_METRIC
		}
	}
	return nil
}

func initStateOfDb(tx *gorm.DB) error {
	if err := initMetrics(tx); err != nil {
		return fmt.Errorf("initMetrics: %w", err)
	}
	return nil
}
