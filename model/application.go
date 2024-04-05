package model

import "time"

type Apps struct {
	App
	GitLabJenkins
}

type GitLabJenkins struct {
	HasGitLab  bool   `json:"has_gitlab"`
	HasJenkins bool   `json:"has_jenkins"`
	Visibility string `json:"visibility"`
}

type App struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	AppName     string `json:"app_name"`
	RepoName    string `json:"repo_name" gorm:"column:repo_name"`
	Lang        string `json:"lang" gorm:"column:lang"`
	Type        string `json:"type"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
}

func (*App) TableName() string {
	return "application"
}
