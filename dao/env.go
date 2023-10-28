package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"kubea-demo/db"
	"kubea-demo/model"
)

var Env env

type env struct{}

type Envs struct {
	Items []*model.Env `json:"items"`
	Total int          `json:"total"`
}

// List 列表
// envName用于模糊查询，过滤
// page，limit用于分页
func (*env) List(envName string, page, limit int) (*Envs, error) {
	//计算分页
	startSet := (page - 1) * limit

	//定义返回值的内容
	var (
		envList = make([]*model.Env, 0)
		total   = 0
	)

	//数据库查询，先查total
	tx := db.GORM.Model(model.Env{}).
		Where("name like ? ", "%"+envName+"%").
		Count(&total)

	if tx.Error != nil {
		logger.Error("获取Env列表失败," + tx.Error.Error())
		return nil, errors.New("获取Env列表失败," + tx.Error.Error())
	}

	//数据库查询，再查数据
	//当limit=10， total一定是10，因为count会在过滤和分页后执行
	tx = db.GORM.Model(model.Env{}).
		Where("name like ?", "%"+envName+"%").
		Limit(limit).
		Offset(startSet).
		Order("name").
		Find(&envList)
	if tx.Error != nil {
		logger.Error("获取Env列表失败," + tx.Error.Error())
		return nil, errors.New("获取Env列表失败," + tx.Error.Error())
	}

	return &Envs{
		Items: envList,
		Total: total,
	}, nil
}

//// Get 查询单个
//func (*env) Get(envName uint) (*model.Env, bool, error) {
//	data := new(model.Env)
//	tx := db.GORM.Where("name = ?", envName).First(&data)
//	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
//		return nil, false, nil
//	}
//
//	if tx.Error != nil {
//		logger.Error("查询Env信息失败," + tx.Error.Error())
//		return nil, false, errors.New("查询Env信息失败," + tx.Error.Error())
//	}
//
//	return data, true, nil
//}

// Has 根据环境名查询，用于代码层去重，查询账号信息
func (*env) Has(envName string) (*model.Env, bool, error) {
	data := new(model.Env)
	tx := db.GORM.Where("name = ?", envName).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据名称查询Env失败," + tx.Error.Error())
		return nil, false, errors.New("根据名称查询Env失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Add 新增
func (*env) Add(e *model.Env) error {
	tx := db.GORM.Create(&e)
	if tx.Error != nil {
		logger.Error("新增Env信息失败," + tx.Error.Error())
		return errors.New("新增Env信息失败," + tx.Error.Error())
	}

	return nil
}

// Update 更新
func (*env) Update(e *model.Env) error {
	tx := db.GORM.Model(&model.Env{}).Where("id = ?", e.ID).Updates(&e)
	if tx.Error != nil {
		logger.Error("更新Env信息失败," + tx.Error.Error())
		return errors.New("更新Env信息失败," + tx.Error.Error())
	}

	return nil
}

// Delete 删除
func (*env) Delete(id uint) error {
	data := new(model.Env)
	data.ID = id
	tx := db.GORM.Delete(&data)
	if tx.Error != nil {
		logger.Error("删除Env信息失败," + tx.Error.Error())
		return errors.New("删除Env信息失败," + tx.Error.Error())
	}

	return nil
}
