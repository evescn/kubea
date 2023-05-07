package model

import "time"

type Chart struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Name     string `json:"name"`
	FileName string `json:"file_name" gorm:"column: file_name"`
	IconUrl  string `json:"icon_url" gorm:"column: icon_url"`
	Version  string `json:"version"`
	Describe string `json:"describe"`
}

// TableName 自定义表名
func (*Chart) TableName() string {
	return "helm_chart"
}
