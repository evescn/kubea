package model

import "time"

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	UserName string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Role     uint   `json:"role"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (*User) TableName() string {
	return "user"
}
