package model

import "time"

type Role struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	RoleName    string `json:"role_name" gorm:"unique;not null" gorm:"column:role_name"`
	Description string `json:"description"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (*Role) TableName() string {
	return "role"
}
