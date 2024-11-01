package db

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Name        string
	Comment     string
	MetricType  uint8
	MetricValue uint64
	//
	CategoryID uint
	Category   *Category
	//
	Proteins uint64
	Carbs    uint64
	Fats     uint64
	//
	Prices         []ItemPrice
	CurrentPriceID uint
	CurrentPrice   *ItemPrice
}

func GetItem(id uint, preloadPrices bool) (*Item, error) {
	if udb == nil {
		return nil, ERROR_DB_NOT_INIT
	}
	db := udb
	if preloadPrices {
		db = db.Preload("Prices")
	}
	var item Item
	err := db.Preload("Category").Preload("CurrentPrice").First(&item, id).Error
	return &item, err
}

func GetItemToRootCat(id uint, preloadPrices bool) (*Item, error) {
	if udb == nil {
		return nil, ERROR_DB_NOT_INIT
	}
	db := udb
	if preloadPrices {
		db = db.Preload("Prices")
	}
	var item Item
	err := db.Preload("Category.Parent", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Parent", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Parent")
		})
	}).Preload("CurrentPrice").First(&item, id).Error
	return &item, err
}
