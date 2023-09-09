package service

import (
	"errors"
	"kubea-demo/dao"
	"kubea-demo/model"
)

var App app

type app struct{}

// List 列表
func (*app) List(appName string, page, limit int) (*dao.Apps, error) {
	return dao.App.List(appName, page, limit)
}

// GetAll 查询所有应用
func (*app) GetAll() ([]*model.App, error) {
	return dao.App.GetAll()
}

// GetRepo 查询所有仓库信息
func (*app) GetRepo() ([]string, error) {
	return dao.App.GetRepo()
}

// GetApp 根据仓库信息查询App信息
func (*app) GetApp(repoName string) ([]string, error) {
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
