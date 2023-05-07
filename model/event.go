package model

import (
	"time"
)

type Event struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Name      string     `json:"name"`
	Kind      string     `json:"kind"`
	Namespace string     `json:"namespace"`
	Rtype     string     `json:"rtype"`
	Reason    string     `json:"reason"`
	Message   string     `json:"message"`
	EventTime *time.Time `json:"event_time"`
	Cluster   string     `json:"cluster"`
}

// TableName 自定义表名
func (*Event) TableName() string {
	return "k8s_event"
}
