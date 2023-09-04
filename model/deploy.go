package model

import (
	"time"
)

type Deploy struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	En        string `json:"en"`
	Tag       string `json:"tag"`
	Status    int    `json:"status"`
	Deployer  string `json:"deployer"`
	StartTime string `json:"start_time" gorm:"column:start_time"`
	Duration  string `json:"duration"`
	//应用与发布数据是一对多关系
	AppId     uint   `json:"app_id" gorm:"column:app_id"`
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
