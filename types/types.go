package types

import "time"

// User struct for requests
type User struct {
	Username string `json:"username" binding:"required" example:"testUser"`
	Password string `json:"password" binding:"required" example:"strongPassLol"`
}

// User Account
type Account struct {
	ID       uint   `json:"id" example:"1"`
	Token    string `json:"token" example:"Fvs-MnxiEs5dnqMp2mSDIJigPbiIUs6Snk1xxiqPmUc"`
	Username string `json:"username" example:"testUser"`
}

type Message struct {
	Info string `json:"info" example:"Success!"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Error: you stink"`
}

type DbCard struct {
	ID             uint   `json:"id" example:"1"`
	Name           string `json:"name" example:"CreditCard"`
	Balance        uint64 `json:"balance" example:"1000"`
	HaveCreditLine bool   `json:"have_credit_line" example:"true"`
	CreditLine     uint64 `json:"credit_line" example:"500000"`
}

type DbCategory struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"Moldova"`
	// Parent is used as a infinite sub-category structure
	//  Can be 0
	ParentID uint `json:"parent_id" example:"0"`
}

type DbDebt struct {
	ID       uint      `json:"id" example:"1"`
	CardID   uint      `json:"card_id" example:"1"`
	Comment  string    `json:"comment" example:"pizza"`
	Value    uint64    `json:"value" example:"20000"`
	IOwe     bool      `json:"i_owe" example:"true"`
	Date     time.Time `json:"date" example:"29/11/2001 12:00"`
	DateEnd  time.Time `json:"date_end" example:"29/12/2001 12:00"`
	Finished bool      `json:"finished" example:"false"`
}

type DbIncome struct {
	ID      uint      `json:"id" example:"1"`
	CardID  uint      `json:"card_id" example:"1"`
	Comment string    `json:"comment" example:"pizza"`
	Value   uint64    `json:"value" example:"20000"`
	Date    time.Time `json:"date" example:"29/11/2001 12:00"`
	UserID  uint      `json:"user_id" example:"1"`
}

type DbType struct {
	ID      uint   `json:"id" example:"1"`
	Name    string `json:"name" example:"Medicine"`
	Comment string `json:"comment" example:""`
	Color   string `json:"color" example:"red"`
}

type DbPayment struct {
	ID         uint      `json:"id" example:"1"`
	CardID     uint      `json:"card_id" example:"1"`
	CategoryID uint      `json:"category_id" example:"1"`
	Title      string    `json:"title" example:"Veggies"`
	Descr      string    `json:"description" example:""`
	Note       string    `json:"not" example:"I'm a teapot"`
	Date       time.Time `json:"date" example:"29/11/2001 12:00"`
}

type Session struct {
	ID     string `json:"id"`
	UserID uint   `json:"user_id" example:"1"`
}

type DbExpense struct {
	ID      uint      `json:"id" example:"1"`
	CardID  uint      `json:"card_id" example:"1"`
	TypeID  uint      `json:"type_id" example:"1"`
	Value   uint64    `json:"value" example:"20000"`
	Comment string    `json:"comment" example:"pizza"`
	Date    time.Time `json:"date" example:"29/11/2001 12:00"`
}

type DbTransfer struct {
	ID         uint      `json:"id" example:"1"`
	FromCardID uint      `json:"from_card_id" example:"1"`
	ToCardID   uint      `json:"to_card_id" example:"1"`
	Value      uint64    `json:"value" example:"20000"`
	Date       time.Time `json:"date" example:"29/11/2001 12:00"`
}
