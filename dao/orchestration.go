package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"kubea-demo/db"
	"kubea-demo/model"
)

var Orchestration orchestration

type orchestration struct{}

// 业务使用的，返回的结构体
// 因为编排管理需要app_name的字段，而我们的数据库种没有这个字段
type OrchestrationRes struct {
	model.Orchestration
	AppName string `json:"app_name"`
}

// 列表专用
type Orchestrations struct {
	Items []*OrchestrationRes `json:"items"`
	Total int                 `json:"total"`
}

// List 列表
func (*orchestration) List(en, appName string, page, limit int) (*Orchestrations, error) {
	startSet := (page - 1) * limit

	var (
		orchestrationList = make([]*OrchestrationRes, 0)
		total             = 0
	)

	//先 count
	query := db.GORM.Model(&model.Orchestration{}).
		Select("orchestration.*, application.app_name").
		Joins("left join application on orchestration.app_id = application.id").
		Where("orchestration.en like ?", "%"+en+"%")
	if appName != "" {
		query = query.Where("application.app_name like ?", "%"+appName+"%")
	}
	tx := query.Count(&total)
	if tx.Error != nil {
		logger.Error("获取Orchestration列表失败," + tx.Error.Error())
		return nil, errors.New("获取Orchestration列表失败," + tx.Error.Error())
	}

	//分页数据
	tx = query.Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&orchestrationList)
	if tx.Error != nil {
		logger.Error("获取Orchestration列表失败," + tx.Error.Error())
		return nil, errors.New("获取Orchestration列表失败," + tx.Error.Error())
	}

	return &Orchestrations{
		Items: orchestrationList,
		Total: total,
	}, nil
}

// Has 查询单个
func (*orchestration) Has(en string, appId uint) (*model.Orchestration, bool, error) {
	data := new(model.Orchestration)
	tx := db.GORM.Where("app_id = ? and en = ?", appId, en).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("查询Orchestration失败," + tx.Error.Error())
		return nil, false, errors.New("查询Orchestration失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Update 更新
func (*orchestration) Update(orch *model.Orchestration) error {
	tx := db.GORM.Model(&model.Orchestration{}).Where("id = ?", orch.ID).Updates(&orch)
	if tx.Error != nil {
		logger.Error("更新Orchestration失败," + tx.Error.Error())
		return errors.New("更新Orchestration失败," + tx.Error.Error())
	}

	return nil
}

// Add 新增
func (*orchestration) Add(orch *model.Orchestration) error {
	tx := db.GORM.Create(&orch)
	if tx.Error != nil {
		logger.Error("新增Orchestration失败," + tx.Error.Error())
		return errors.New("新增Orchestration失败," + tx.Error.Error())
	}

	return nil
}

// Delete 删除
func (*orchestration) Delete(id uint) error {
	data := new(model.Orchestration)
	data.ID = id
	tx := db.GORM.Delete(&data)
	if tx.Error != nil {
		logger.Error("删除Orchestration失败," + tx.Error.Error())
		return errors.New("删除Orchestration失败," + tx.Error.Error())
	}

	return nil
}
