package service

import (
	"errors"
	"kubea/dao"
	"kubea/model"
)

var App app

type app struct{}

// List 列表
func (*app) List(appName, repoName string, page, limit int) (*dao.Apps, error) {
	return dao.App.List(appName, repoName, page, limit)
}

// Get 查询单个应用信息
func (*app) Get(repoName, appName string) (*model.App, error) {
	data, has, err := dao.App.Has(repoName, appName)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("查询无此应用")
	}
	return data, nil
}

// GetAll 查询所有应用
func (*app) GetAll() ([]*model.App, error) {
	return dao.App.GetAll()
}

// GetApp 根据仓库信息查询App信息
func (*app) GetApp(repoName string) ([]*model.App, error) {
	data, has, err := dao.App.GetApp(repoName)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("repo 仓库下没有数据")
	}
	return data, nil
}

// Update 更新
func (*app) Update(app *model.App) error {
	return dao.App.Update(app)
}

// Add 新增
func (*app) Add(app *model.App) error {
	//查看应用名是否存在
	_, has, err := dao.App.Has(app.RepoName, app.AppName)
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
