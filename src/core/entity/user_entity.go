package entity

type UserEntity struct {
	BaseEntity
	FullName string
	Email    string
	Password string
}
