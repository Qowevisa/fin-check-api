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

type ErrorResponse struct {
	Message string `json:"message" example:"Error: you stink"`
}
