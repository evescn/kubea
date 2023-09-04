package model

import (
	"time"
)

type Build struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	RepoName string `json:"repo_name"`
	Tag      string `json:"tag"`
	Builder  string `json:"builder"`
	BuildUrl string `json:"build_url" gorm:"column:build_url"`
}

func (*Build) TableName() string {
	return "build"
}
