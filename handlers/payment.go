package handlers

import (
	"fmt"
	"log"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var itemBoughtTransform func(inp *db.ItemBought) types.ItemBought = func(inp *db.ItemBought) types.ItemBought {
	return types.ItemBought{}
}

var paymentTransform func(inp *db.Payment) types.Payment = func(inp *db.Payment) types.Payment {
	var items []types.ItemBought
	for _, item := range inp.Items {
		items = append(items, itemBoughtTransform(&item))
	}
	return types.Payment{
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
// @Param payment body types.Payment true "Payment"
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

	var updates types.Payment
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
				if err := dbc.Debug().Delete(deleteIt).Error; err != nil {
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
	for _, uItemBought := range updates.Items {
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
			newItemBought := &db.ItemBought{
				ItemID:      newItem.ID,
				PaymentID:   payment.ID,
				TypeID:      uItemBought.TypeID,
				Quantity:    uItemBought.Quantity,
				TotalCost:   newItem.Price * uint64(uItemBought.Quantity),
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
			// TODO: check if Item has same userID and potentially update Item
		}
	}

	c.JSON(200, types.Message{Info: fmt.Sprintf("Entity with %d ID is created successfully!", payment.ID)})
}
