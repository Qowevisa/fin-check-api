package db

type Identifiable interface {
	GetID() uint
	SetID(id uint)
}

type UserIdentifiable interface {
	GetID() uint
	GetUserID() uint
	SetUserID(userID uint)
}
