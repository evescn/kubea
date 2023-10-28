package model

import "time"

type RoleMenuRelation struct {
	ID           uint `json:"id" gorm:"primary_key"`
	RoleID       uint `json:"role_id" gorm:"column:role_id"`
	PageID       uint `json:"page_id" gorm:"column:page_id"`
	SubPageID    uint `json:"sub_page_id" gorm:"column:sub_page_id"`
	SubSubPageID uint `json:"sub_sub_page_id" gorm:"column:sub_sub_page_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (*RoleMenuRelation) TableName() string {
	return "role_menu_relation"
}
