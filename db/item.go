package db

import (
	"errors"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name        string
	Comment     string
	Price       uint64
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
	Prices         []ItemPrice `gorm:"constraint:OnDelete:CASCADE;"`
	CurrentPriceID uint
	CurrentPrice   *ItemPrice
	//
	TypeID uint
	Type   *Type
	UserID uint
	User   *User
}

func (i Item) __internalBelogingToPayment() {}

// Implements db.UserIdentifiable:1
func (i Item) GetID() uint {
	return i.ID
}

// Implements db.UserIdentifiable:2
func (i Item) GetUserID() uint {
	return i.UserID
}

// Implements db.UserIdentifiable:3
func (i *Item) SetUserID(id uint) {
	i.UserID = id
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

var (
	ERROR_ITEM_PRICE_ISZERO     = errors.New("Item's Price is zero")
	ERROR_ITEM_ITEMPRICE_INTERR = errors.New("Item's ItemPrice ID is zero after creating")
)

func (i *Item) BeforeCreate(tx *gorm.DB) error {
	if i.Price == 0 {
		return ERROR_ITEM_PRICE_ISZERO
	}
	return nil
}

func (i *Item) AfterCreate(tx *gorm.DB) error {
	if i.CurrentPriceID == 0 {
		itemPrice := &ItemPrice{
			ItemID: i.ID,
			Price:  i.Price,
		}
		if err := tx.Create(itemPrice).Error; err != nil {
			return err
		}
		if itemPrice.ID == 0 {
			return ERROR_ITEM_ITEMPRICE_INTERR
		}
		i.CurrentPriceID = itemPrice.ID
		if err := tx.Save(i).Error; err != nil {
			return err
		}
	}
	return nil
}
