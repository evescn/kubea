package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"kubea/db"
	"kubea/model"
)

var App app

type app struct{}

// 一般情况下，对于列表接口我们要返回total和items,total是过滤后，分页前的数量
type Apps struct {
	Items []*model.App `json:"items"`
	Total int          `json:"total"`
}

// List 列表
// appName用于模糊查询，过滤
// page，limit用于分页
// 默认desc倒序
func (*app) List(appName, repoName string, page, limit int) (*Apps, error) {
	//计算分页
	startSet := (page - 1) * limit

	//定义返回值的内容
	var (
		appList = make([]*model.App, 0)
		total   = 0
	)

	//数据库查询，先查total
	tx := db.GORM.Model(model.App{}).
		Where("repo_name like ?", "%"+repoName+"%").
		Where("app_name like ? ", "%"+appName+"%").
		Count(&total)

	if tx.Error != nil {
		logger.Error("获取Application列表失败," + tx.Error.Error())
		return nil, errors.New("获取Application列表失败," + tx.Error.Error())
	}

	//数据库查询，再查数据
	//当limit=10， total一定是10，因为count会在过滤和分页后执行
	tx = db.GORM.Model(model.App{}).
		Where("repo_name like ?", "%"+repoName+"%").
		Where("app_name like ?", "%"+appName+"%").
		Limit(limit).
		Offset(startSet).
		Order("app_name").
		Find(&appList)
	if tx.Error != nil {
		logger.Error("获取Application列表失败," + tx.Error.Error())
		return nil, errors.New("获取Application列表失败," + tx.Error.Error())
	}

	return &Apps{
		Items: appList,
		Total: total,
	}, nil
}

// Get 查询单个
func (*app) Get(appId uint) (*model.App, bool, error) {
	data := new(model.App)
	tx := db.GORM.Where("id = ?", appId).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("查询Application失败," + tx.Error.Error())
		return nil, false, errors.New("查询Application失败," + tx.Error.Error())
	}

	return data, true, nil
}

// GetAll 查询所有应用
func (*app) GetAll() ([]*model.App, error) {
	data := make([]*model.App, 0)
	tx := db.GORM.Find(&data)
	if tx.Error != nil {
		logger.Error("查询所有Application失败," + tx.Error.Error())
		return nil, errors.New("查询所有Application失败," + tx.Error.Error())
	}

	return data, nil
}

// GetApp 根据仓库名查询所有App
func (*app) GetApp(repo string) ([]*model.App, bool, error) {
	data := make([]*model.App, 0)
	tx := db.GORM.Model(&model.App{}).
		Where("repo_name = ?", repo).
		Find(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据仓库名查询Application失败," + tx.Error.Error())
		return nil, false, errors.New("根据仓库名查询Application失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Has 根据应用名查询，用于代码层去重
func (*app) Has(repoName, appName string) (*model.App, bool, error) {
	data := new(model.App)
	tx := db.GORM.Where("repo_name = ? and app_name = ?", repoName, appName).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据应用名查询Application失败," + tx.Error.Error())
		return nil, false, errors.New("根据应用名查询Application失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Update 更新
func (*app) Update(app *model.App) error {
	tx := db.GORM.Model(&model.App{}).Where("id = ?", app.ID).Updates(&app)
	if tx.Error != nil {
		logger.Error("更新Application失败," + tx.Error.Error())
		return errors.New("更新Application失败," + tx.Error.Error())
	}

	return nil
}

// Add 新增
func (*app) Add(app *model.App) error {
	tx := db.GORM.Create(&app)
	if tx.Error != nil {
		logger.Error("新增Application失败," + tx.Error.Error())
		return errors.New("新增Application失败," + tx.Error.Error())
	}

	return nil
}

// Delete 删除
func (*app) Delete(id uint) error {
	data := new(model.App)
	data.ID = id
	tx := db.GORM.Delete(&data)
	if tx.Error != nil {
		logger.Error("删除Application失败," + tx.Error.Error())
		return errors.New("删除Application失败," + tx.Error.Error())
	}

	return nil
}
