package handlers

import (
	"fmt"
	"log"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var itemBoughtTransform func(inp *db.ItemBought) types.DbItemBought = func(inp *db.ItemBought) types.DbItemBought {
	var item types.DbItem
	var price uint64 = 0
	if inp.Item != nil {
		item = itemTransform(inp.Item)
		price = inp.Item.Price
	}
	return types.DbItemBought{
		ID:          inp.ID,
		ItemID:      inp.ItemID,
		PaymentID:   inp.PaymentID,
		TypeID:      inp.TypeID,
		Price:       price,
		Quantity:    inp.Quantity,
		TotalCost:   inp.TotalCost,
		MetricType:  inp.MetricType,
		MetricValue: inp.MetricValue,
		Item:        item,
	}
}

var paymentTransform func(inp *db.Payment) types.DbPayment = func(inp *db.Payment) types.DbPayment {
	var items []types.DbItemBought
	for _, item := range inp.Items {
		items = append(items, itemBoughtTransform(&item))
	}
	return types.DbPayment{
		ID:          inp.ID,
		CardID:      inp.CardID,
		CategoryID:  inp.CategoryID,
		Title:       inp.Title,
		Description: inp.Descr,
		Note:        inp.Note,
		Date:        inp.Date,
		Items:       items,
	}
}

// @Summary Add payment
// @Description Add payment
// @Tags payment
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param payment body types.DbPayment true "Payment"
// @Success 200 {object} types.Message
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /payment/add [post]
func PaymentAdd(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var updates types.DbPayment
	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Printf("err is %v\n", err)
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}
	// As this handler will likely create more than one row in database we need to
	// create some sort of defer func that will rollback all created rows
	weNeedRollback := false
	var deletableIfRollback []db.PaymentGroup
	defer func() {
		if weNeedRollback {
			dbc := db.Connect()
			for _, deleteIt := range deletableIfRollback {
				if err := dbc.Delete(deleteIt).Error; err != nil {
					log.Printf("ERROR: dbc.Delete: %v\n", err)
					continue
				}
			}
		}
	}()
	dbc := db.Connect()

	payment := &db.Payment{
		CardID:     updates.CardID,
		CategoryID: updates.CategoryID,
		UserID:     userID,
		Title:      updates.Title,
		Descr:      updates.Description,
		Note:       updates.Note,
		Date:       updates.Date,
	}
	if err := dbc.Debug().Create(payment).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		weNeedRollback = true
		return
	}
	if payment.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "Internal error: ERR.P.A.1"})
		weNeedRollback = true
		return
	}
	deletableIfRollback = append(deletableIfRollback, payment)
	var substractFromCard uint64 = 0
	for _, uItemBought := range updates.Items {
		var totalCost uint64 = 0
		// Creating item and adding it to rollback if itemID is set to 0
		if uItemBought.ItemID == 0 {
			newItem := &db.Item{
				Name:        uItemBought.NewName,
				Comment:     uItemBought.NewComment,
				Price:       uItemBought.Price,
				MetricType:  uItemBought.MetricType,
				MetricValue: uItemBought.MetricValue,
				CategoryID:  updates.CategoryID,
				TypeID:      uItemBought.TypeID,
				UserID:      userID,
			}
			if err := dbc.Create(newItem).Error; err != nil {
				c.JSON(500, types.ErrorResponse{Message: err.Error()})
				weNeedRollback = true
				return
			}
			if newItem.ID == 0 {
				c.JSON(500, types.ErrorResponse{Message: "Internal error: ERR.P.A.2"})
				weNeedRollback = true
				return
			}
			deletableIfRollback = append(deletableIfRollback, newItem)
			totalCost = newItem.Price * uint64(uItemBought.Quantity)
			newItemBought := &db.ItemBought{
				ItemID:      newItem.ID,
				PaymentID:   payment.ID,
				TypeID:      uItemBought.TypeID,
				Quantity:    uItemBought.Quantity,
				TotalCost:   totalCost,
				MetricType:  uItemBought.MetricType,
				MetricValue: uItemBought.MetricValue,
			}
			if err := dbc.Create(newItemBought).Error; err != nil {
				c.JSON(500, types.ErrorResponse{Message: err.Error()})
				weNeedRollback = true
				return
			}
			if newItemBought.ID == 0 {
				c.JSON(500, types.ErrorResponse{Message: "Internal error: ERR.P.A.3"})
				weNeedRollback = true
				return
			}
			deletableIfRollback = append(deletableIfRollback, newItemBought)
			newItemBought.Item = newItem
		} else {
			dbItem := &db.Item{}
			if err := dbc.Find(dbItem, uItemBought.ItemID).Error; err != nil {
				c.JSON(500, types.ErrorResponse{Message: err.Error()})
				weNeedRollback = true
				return
			}
			if dbItem.UserID != userID {
				c.JSON(500, types.ErrorResponse{Message: "ItemID is not your"})
				weNeedRollback = true
				return
			}
			if dbItem.Price != uItemBought.Price {
				dbItem.Price = uItemBought.Price
				if err := dbc.Save(dbItem).Error; err != nil {
					c.JSON(500, types.ErrorResponse{Message: err.Error()})
					weNeedRollback = true
					return
				}
			}
			totalCost = dbItem.Price * uint64(uItemBought.Quantity)
			newItemBought := &db.ItemBought{
				ItemID:      dbItem.ID,
				PaymentID:   payment.ID,
				TypeID:      uItemBought.TypeID,
				Quantity:    uItemBought.Quantity,
				TotalCost:   totalCost,
				MetricType:  uItemBought.MetricType,
				MetricValue: uItemBought.MetricValue,
			}
			if err := dbc.Create(newItemBought).Error; err != nil {
				c.JSON(500, types.ErrorResponse{Message: err.Error()})
				weNeedRollback = true
				return
			}
			if newItemBought.ID == 0 {
				c.JSON(500, types.ErrorResponse{Message: "Internal error: ERR.P.A.4"})
				weNeedRollback = true
				return
			}
			deletableIfRollback = append(deletableIfRollback, newItemBought)
			newItemBought.Item = dbItem
		}
		// As totalCost is calculated either way AND db.Item.BeforeSave have db.Item.Price check
		//  we can gently assume that totalCost != 0
		substractFromCard += totalCost
	}

	// Checks for db.Card.UserID != db.Payment.UserID for payment is done in db.Payment hooks
	card := &db.Card{}
	if err := dbc.Find(card, payment.CardID).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		weNeedRollback = true
		return
	}
	// TODO: Examine if that would be better to calculate substractFromCard
	//  BEFORE update.Items loop and therefore this check can be done before
	//  processing payment's ItemBought entities
	if card.Balance < substractFromCard {
		c.JSON(500, types.ErrorResponse{Message: fmt.Sprintf("Card %s (%s) has Balance lower than overall Payment cost. Rollback.", card.Name, card.LastDigits)})
		weNeedRollback = true
		return
	}
	card.Balance -= substractFromCard
	if err := dbc.Save(card).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		weNeedRollback = true
		return
	}

	c.JSON(200, types.Message{Info: fmt.Sprintf("Entity with %d ID is created successfully!", payment.ID)})
}

// @Summary Get all payments for user
// @Description Get all payments for user
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.DbPayment
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /payment/all [get]
func PaymentGetAll(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	dbc := db.Connect()
	var entities []*db.Payment
	if err := dbc.Preload("Items.Item").Find(&entities, db.Payment{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.DbPayment
	for _, entity := range entities {
		ret = append(ret, paymentTransform(entity))
	}
	c.JSON(200, ret)
}
