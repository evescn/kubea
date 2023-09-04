package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"kubea-demo/db"
	"kubea-demo/model"
)

var Build build

type build struct{}

// 一般情况下，对于列表接口我们要返回total和items,total是过滤后，分页前的数量
type Builds struct {
	Items []*model.Build `json:"items"`
	Total int            `json:"total"`
}

// List 列表
func (*build) List(repoName string, page, limit int) (*Builds, error) {
	startSet := (page - 1) * limit

	var (
		buildList = make([]*model.Build, 0)
		total     = 0
	)

	tx := db.GORM.Model(&model.Build{}).Where("repo_name like ?", "%"+repoName+"%").Count(&total)
	if tx.Error != nil {
		logger.Error("获取Build列表失败," + tx.Error.Error())
		return nil, errors.New("获取Build列表失败," + tx.Error.Error())
	}

	tx = db.GORM.Model(&model.Build{}).Where("repo_name like ?", "%"+repoName+"%").Limit(limit).Offset(startSet).Order("id desc").Find(&buildList)
	if tx.Error != nil {
		logger.Error("获取Build列表失败," + tx.Error.Error())
		return nil, errors.New("获取Build列表失败," + tx.Error.Error())
	}

	return &Builds{
		Items: buildList,
		Total: total,
	}, nil
}

// GetAllTagByRepo 根据repo查询所有tag
func (*build) GetAllTagByRepo(repoName string) ([]string, error) {
	data := make([]*model.Build, 0)
	tx := db.GORM.Model(&model.Build{}).Where("repo_name = ?", "%"+repoName+"%").Order("id desc").Find(&data)
	if tx.Error != nil {
		logger.Error("根据repo获取tag失败," + tx.Error.Error())
		return nil, errors.New("根据repo获取tag失败," + tx.Error.Error())
	}

	mp := make(map[string]int, 0)
	list := make([]string, 0)
	for _, val := range data {
		//如果这个map的值为1，则是重复的，直接跳过
		if _, ok := mp[val.Tag]; ok {
			continue
		}
		//如果把tag加入到list中，则做个标记，map的值为1
		list = append(list, val.Tag)
		mp[val.Tag] = 1
	}

	return list, nil
}

// Has 查询单个
func (*build) Has(repoName, tag string) (*model.Build, bool, error) {
	data := new(model.Build)
	tx := db.GORM.Where("repo_name = ? and tag = ?", repoName, tag).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("查询Build失败," + tx.Error.Error())
		return nil, false, errors.New("查询Build失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Add 新增
func (*build) Add(b *model.Build) error {
	tx := db.GORM.Create(&b)
	if tx.Error != nil {
		logger.Error("新增Build失败," + tx.Error.Error())
		return errors.New("新增Build失败," + tx.Error.Error())
	}

	return nil
}

// Update 更新
func (*build) Update(b *model.Build) error {
	tx := db.GORM.Model(&model.Build{}).Where("id = ?", b.ID).Updates(&b)
	if tx.Error != nil {
		logger.Error("更新Build失败," + tx.Error.Error())
		return errors.New("更新Build失败," + tx.Error.Error())
	}

	return nil
}
