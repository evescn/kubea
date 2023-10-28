package model

import "time"

type Menu struct {
	ID          uint   `json:"page_id" gorm:"primary_key"`
	Path        string `json:"page_path" gorm:"unique;not null"`
	Name        string `json:"page_name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// TableName 自定义表名
func (*Menu) TableName() string {
	return "menu"
}
