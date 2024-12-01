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
	gormDB, err := gorm.Open(sqlite.Open("fin-check.db"), &gorm.Config{
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
	gormDB.AutoMigrate(&ExchangeRate{})
	gormDB.AutoMigrate(&SettingsTypeFilter{})
	return newUDB
}

var (
	CANT_FIND_METRIC   = errors.New("Can't find proper metrics in database")
	CANT_FIND_CURRENCY = errors.New("Can't find proper currencies in database")
)

func checkSeededValues[T Identifiable](whatToCheck []*T, errorIfNotFound error, tx *gorm.DB) error {
	var valuesInDB []T
	if err := tx.Find(&valuesInDB).Error; err != nil {
		return err
	}
	if len(valuesInDB) == 0 {
		for _, v := range whatToCheck {
			if err := tx.Create(v).Error; err != nil {
				return err
			}
		}
		return nil
	}
	for _, v := range whatToCheck {
		var tmp T
		if err := tx.Find(&tmp, v).Error; err != nil {
			return err
		}
		if tmp.GetID() == 0 {
			return errorIfNotFound
		}
	}
	return nil
}

func initMetrics(tx *gorm.DB) error {
	metricsThatNeeded := []*Metric{
		&Metric{Name: "None", Short: "pcs", Type: METRIC_TYPE_NONE},
		&Metric{Name: "Gram", Short: "g", Type: METRIC_TYPE_GRAM},
		&Metric{Name: "Kilogram", Short: "kg", Type: METRIC_TYPE_KILOGRAM},
		&Metric{Name: "Liter", Short: "l", Type: METRIC_TYPE_LITER},
	}
	return checkSeededValues(metricsThatNeeded, CANT_FIND_METRIC, tx)
}

func initCurrencies(tx *gorm.DB) error {
	currsThatNeeded := []*Currency{
		{Name: "Dollar", Symbol: "$", ISOName: "USD"},
		{Name: "Moldavian Leu", Symbol: "L", ISOName: "MDL"},
		{Name: "Romanian Leu", Symbol: "RL", ISOName: "RON"},
		{Name: "Polish Zloty", Symbol: "zł", ISOName: "PLN"},
		{Name: "Ukrainian Hryvnia", Symbol: "₴", ISOName: "UAH"},
		{Name: "Euro", Symbol: "€", ISOName: "EUR"},
		{Name: "Russian Ruble", Symbol: "₽", ISOName: "RUB"},
		{Name: "Kazakhstani Tenge", Symbol: "₸", ISOName: "KZT"},
		{Name: "Chinese Yuan", Symbol: "¥", ISOName: "CNY"},
	}
	return checkSeededValues(currsThatNeeded, CANT_FIND_CURRENCY, tx)
}

func initStateOfDb(tx *gorm.DB) error {
	if err := initMetrics(tx); err != nil {
		return fmt.Errorf("initMetrics: %w", err)
	}
	if err := initCurrencies(tx); err != nil {
		return fmt.Errorf("initCurrencies: %w", err)
	}
	return nil
}
