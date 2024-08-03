package types

// User struct for requests
type User struct {
	Username string `json:"username" binding:"required" example:"testUser"`
	Password string `json:"password" binding:"required" example:"strongPassLol"`
}

// User Account
type Account struct {
	ID    uint   `json:"id" example:"1"`
	Token string `json:"token" example:"Fvs-MnxiEs5dnqMp2mSDIJigPbiIUs6Snk1xxiqPmUc="`
}

type Message struct {
	Message string `json:"message" example:"Success!"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Error: you stink"`
}

type DbCard struct {
	ID             uint   `json:"id" example:"1"`
	Name           string `json:"name" example:"CreditCard"`
	Value          uint64 `json:"value" example:"1000"`
	HaveCreditLine bool   `json:"have_credit_line" example:"true"`
	CreditLine     uint64 `json:"credit_line" example:"500000"`
}

type DbCategory struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"Moldova"`
	// Parent is used as a infinite sub-category structure
	//  Can be 0
	ParentID uint `json:"parent_id" example:"0"`
	UserID   uint `json:"user_id" example:"1"`
}
