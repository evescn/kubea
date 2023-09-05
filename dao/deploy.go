package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"kubea-demo/db"
	"kubea-demo/model"
)

var Deploy deploy

type deploy struct{}

type DeployRes struct {
	model.Deploy
	AppName string `json:"app_name"`
}

type Deploys struct {
	Items []*DeployRes `json:"items"`
	Total int          `json:"total"`
}

// List 列表
func (*deploy) List(en, appName string, page, limit int) (*Deploys, error) {
	startSet := (page - 1) * limit

	var (
		deployList = make([]*DeployRes, 0)
		total      = 0
	)

	query := db.GORM.Model(&model.Deploy{}).
		Select("deploy.*, application.app_name").
		Joins("left join application on deploy.app_id = application.id").
		Where("deploy.en like ?", "%"+en+"%")
	if appName != "" {
		query = query.Where("application.app_name like ?", "%"+appName+"%")
	}
	tx := query.Count(&total)
	if tx.Error != nil {
		logger.Error("获取Build列表失败," + tx.Error.Error())
		return nil, errors.New("获取Build列表失败," + tx.Error.Error())
	}

	//分页数据
	tx = query.Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&deployList)
	if tx.Error != nil {
		logger.Error("获取Build列表失败," + tx.Error.Error())
		return nil, errors.New("获取Build列表失败," + tx.Error.Error())
	}

	return &Deploys{
		Items: deployList,
		Total: total,
	}, nil

}

// Get 查询单个
func (*deploy) Get(deployId uint) (*model.Deploy, bool, error) {
	data := new(model.Deploy)
	tx := db.GORM.Model(&model.Deploy{}).Where("id = ?", deployId).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if tx.Error != nil {
		logger.Error("查询Deploy失败," + tx.Error.Error())
		return nil, false, errors.New("查询Deploy失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Add 新增
func (*deploy) Add(d *model.Deploy) error {
	tx := db.GORM.Create(&d)
	if tx.Error != nil {
		logger.Error("新增Deploy失败," + tx.Error.Error())
		return errors.New("新增Deploy失败," + tx.Error.Error())
	}

	return nil
}

// Update 更新
func (*deploy) Update(d *model.Deploy) error {
	tx := db.GORM.Model(&model.Deploy{}).Where("id = ?", d.ID).Updates(&d)
	if tx.Error != nil {
		logger.Error("更新Deploy失败," + tx.Error.Error())
		return errors.New("更新Deploy失败," + tx.Error.Error())
	}

	return nil
}

// Delete 删除
func (*deploy) Delete(deployId uint) error {
	data := new(model.Deploy)
	data.ID = deployId
	tx := db.GORM.Delete(&data)
	if tx.Error != nil {
		logger.Error("删除Deploy失败," + tx.Error.Error())
		return errors.New("删除Deploy失败," + tx.Error.Error())
	}

	return nil
}

// GetLog 查询日志
func (*deploy) GetLog(deployId uint) (*model.DeployLog, bool, error) {
	data := new(model.DeployLog)
	tx := db.GORM.Model(&model.Deploy{}).Where("deploy_id = ?", deployId).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if tx.Error != nil {
		logger.Error("查询DeployLog失败," + tx.Error.Error())
		return nil, false, errors.New("查询DeployLog失败," + tx.Error.Error())
	}

	return data, true, nil
}

// UpdateLog 修改日志，追加
func (*deploy) UpdateLog(deployLog *model.DeployLog) error {
	tx := db.GORM.Model(&model.DeployLog{}).Where("deploy_id = ?", deployLog.DeployId).Updates(&deployLog)
	if tx.Error != nil {
		logger.Error("更新DeployLog失败," + tx.Error.Error())
		return errors.New("更新DeployLog失败," + tx.Error.Error())
	}

	return nil
}

// AddLog 新增日志
func (*deploy) AddLog(deployLog *model.DeployLog) error {
	tx := db.GORM.Create(&deployLog)
	if tx.Error != nil {
		logger.Error("新增DeployLog失败," + tx.Error.Error())
		return errors.New("新增DeployLog失败," + tx.Error.Error())
	}

	return nil
}
