package model

import (
	"time"
)

type Deploy struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	En           string `json:"en"`
	Branch       string `json:"branch"`
	Tag          int    `json:"tag"`
	Status       int    `json:"status"`
	StartTime    string `json:"start_time" gorm:"column:start_time"`
	Duration     string `json:"duration"`
	BuildStatus  int    `json:"build_status" gorm:"column:build_status"`
	DeployStatus int    `json:"deploy_status" gorm:"column:deploy_status"`
	Builder      string `json:"builder"`
	BuildUrl     string `json:"build_url" gorm:"column:build_url"`
	//应用与发布数据是一对多关系
	AppId uint `json:"app_id" gorm:"column:app_id"`
}

func (*Deploy) TableName() string {
	return "deploy"
}

type DeployLog struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	//发布日志与发布数据时一对一关系
	DeployId uint   `json:"deploy_id" gorm:"column:deploy_id"`
	Log      string `json:"log"`
}

func (*DeployLog) TableName() string {
	return "deploy_log"
}
