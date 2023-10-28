package model

import "time"

type Password struct {
	ID       uint   `json:"id" gorm:"primary_key" gorm:"column:id"`
	PName    string `json:"p_name" gorm:"column:p_name"`
	Password string `json:"p_password"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// TableName 自定义表名
func (*Password) TableName() string {
	return "password"
}
