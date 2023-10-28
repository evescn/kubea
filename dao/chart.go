package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"kubea/db"
	"kubea/model"
)

var Chart chart

type chart struct{}

type Charts struct {
	Items []*model.Chart `json:"items"`
	Total int            `json:"total"`
}

// GetList 获取列表
func (*chart) GetList(name string, page, limit int) (*Charts, error) {
	//定义分页数据的起始位置
	startSet := (page - 1) * limit
	//定义数据库查询的返回内容
	var (
		chartList = make([]*model.Chart, 0)
		total     = 0
	)

	//数据库查询，Limit方法用于限制条数，Offset方法设置起始位置
	tx := db.GORM.
		Model(&model.Chart{}).
		Where("name like ?", "%"+name+"%").
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&chartList)

	if tx.Error != nil {
		logger.Error(fmt.Sprintf("获取Chart列表失败,%v\n", tx.Error))
		return nil, errors.New(fmt.Sprintf("获取Chart列表失败,%v\n", tx.Error))
	}

	return &Charts{
		Items: chartList,
		Total: total,
	}, nil
}

// Has 查询单个
func (*chart) Has(name string) (*model.Chart, bool, error) {
	data := &model.Chart{}
	tx := db.GORM.Where("name = ?", name).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if tx.Error != nil {
		logger.Error(fmt.Sprintf("查询Chart失败, %v\n", tx.Error))
		return nil, false, errors.New(fmt.Sprintf("查询Chart失败, %v\n", tx.Error))
	}
	return data, true, nil
}

// Add 新增
func (*chart) Add(chart *model.Chart) error {
	tx := db.GORM.Create(&chart)
	if tx.Error != nil {
		logger.Error(fmt.Sprintf("添加Chart失败, %v\n", tx.Error))
		return errors.New(fmt.Sprintf("添加Chart失败, %v\n", tx.Error))
	}
	return nil
}

// Update 更新
func (*chart) Update(chart *model.Chart) error {
	tx := db.GORM.Model(&chart).Updates(&model.Chart{
		Name:     chart.Name,
		FileName: chart.FileName,
		IconUrl:  chart.IconUrl,
		Version:  chart.Version,
		Describe: chart.Describe,
	})
	if tx.Error != nil {
		logger.Error(fmt.Sprintf("更新Chart失败, %v\n", tx.Error))
		return errors.New(fmt.Sprintf("更新Chart失败, %v\n", tx.Error))
	}
	return nil
}

// Delete 删除
func (*chart) Delete(id uint) error {
	data := &model.Chart{}
	data.ID = uint(id)
	tx := db.GORM.Delete(&data)
	if tx.Error != nil {
		logger.Error("删除Chart失败, " + tx.Error.Error())
		return errors.New("删除Chart失败, " + tx.Error.Error())
	}
	return nil
}
