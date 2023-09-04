package model

import "time"

type Orchestration struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	En         string  `json:"en"`
	Replicas   uint    `json:"replicas"`
	LimitMem   string  `json:"limit_mem" gorm:"column:limit_mem"`
	RequestMem string  `json:"request_mem" gorm:"column:request_mem"`
	LimitCpu   float32 `json:"limit_cpu" gorm:"column:limit_cpu"`
	RequestCpu float32 `json:"request_cpu" gorm:"column:request_cpu"`
	Command    string  `json:"command"`
	AppId      uint    `json:"app_id"`
}

func(*Orchestration) TableName() string {
	return "orchestration"
}