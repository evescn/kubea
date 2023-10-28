package model

import "time"

type SubMenu struct {
	ID           uint   `json:"sub_page_id" gorm:"primary_key"`
	Path         string `json:"sub_page_path" gorm:"unique;not null"`
	Name         string `json:"sub_page_name"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	ParentPageID int    `json:"parent_page_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// TableName 自定义表名
func (*SubMenu) TableName() string {
	return "sub_menu"
}
