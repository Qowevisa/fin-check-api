package db

type Identifiable interface {
	GetID() uint
}

type UserIdentifiable interface {
	GetID() uint
	GetUserID() uint
	SetUserID(userID uint)
}

type PaymentGroup interface {
	__internalBelogingToPayment()
}
