package model

import "time"

type SubSubMenu struct {
	ID              uint   `json:"sub_sub_id" gorm:"primary_key"`
	Path            string `json:"sub_sub_page_path" gorm:"unique;not null" `
	Name            string `json:"sub_sub_page_name"`
	Description     string `json:"description"`
	ParentSubPageID int    `json:"parent_sub_page_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// TableName 自定义表名
func (*SubSubMenu) TableName() string {
	return "sub_sub_menu"
}
