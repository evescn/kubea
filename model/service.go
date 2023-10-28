package model

import "time"

type Service struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"svc_name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Eid         uint   `json:"eid"`
	Pid         uint   `json:"pid"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// TableName 自定义表名
func (*Service) TableName() string {
	return "service"
}
