package service

import (
	"errors"
	"kubea-demo/dao"
	"kubea-demo/model"
)

var Orchestration orchestration

type orchestration struct{}

// 列表
func (*orchestration) List(en, appName string, page, limit int) (*dao.Orchestrations, error) {
	return dao.Orchestration.List(en, appName, page, limit)
}

// Update 更新
func (*orchestration) Update(orch *model.Orchestration) error {
	return dao.Orchestration.Update(orch)
}

// Add 新增
func (*orchestration) Add(orch *model.Orchestration) error {
	//查看应用名是否存在
	_, has, err := dao.Orchestration.Has(orch.En, orch.AppId)
	if err != nil {
		return err
	}
	if has {
		return errors.New("该数据已存在，请重新添加")
	}

	//不存在则创建
	if err := dao.Orchestration.Add(orch); err != nil {
		return err
	}

	return nil
}

// Delete 删除
func (*orchestration) Delete(id uint) error {
	return dao.Orchestration.Delete(id)
}
