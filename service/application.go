package service

import (
	"errors"
	"kubea-demo/dao"
	"kubea-demo/model"
)

var App app

type app struct{}

// 列表
func (*app) List(appName string, page, limit int) (*dao.Apps, error) {
	return dao.App.List(appName, page, limit)
}

// GetAll 查询所有应用
func (*app) GetAll() ([]*model.App, error) {
	return dao.App.GetAll()
}

// GetAllTabByApp 查询所有tag
func (*app) GetAllTabByApp(appId uint) ([]string, error) {
	data, has, err := dao.App.Get(appId)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("该数据不存在")
	}

	return dao.Build.GetAllTagByRepo(data.RepoName)
}

// Update 更新
func (*app) Update(app *model.App) error {
	return dao.App.Update(app)
}

// Add 新增
func (*app) Add(app *model.App) error {
	//查看应用名是否存在
	_, has, err := dao.App.Has(app.AppName)
	if err != nil {
		return err
	}
	if has {
		return errors.New("该数据已存在，请重新添加")
	}

	//不存在则创建
	if err := dao.App.Add(app); err != nil {
		return err
	}

	return nil
}

// Delete 删除
func (*app) Delete(id uint) error {
	return dao.App.Delete(id)
}
