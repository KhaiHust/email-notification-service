package model

type UserModel struct {
	BaseModel
	FullName string `gorm:"column:full_name"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

func (UserModel) TableName() string {
	return "users"
}
