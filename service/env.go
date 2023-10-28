package service

import (
	"errors"
	"github.com/wonderivan/logger"
	"kubea-demo/dao"
	"kubea-demo/model"
)

var Env env

type env struct{}

// List 返回环境列表
func (*env) List(envName string, page, limit int) (*dao.Envs, error) {
	return dao.Env.List(envName, page, limit)
}

// Add 创建环境
func (*env) Add(e *model.Env) error {
	// 判断环境是否存在
	_, has, err := dao.Env.Has(e.Name)
	if err != nil {
		return err
	}

	if has {
		logger.Error("当前数据已存在，请重新创建")
		return errors.New("当前数据已存在，请重新创建")
	}

	// 不存在则创建数据
	if err = dao.Env.Add(e); err != nil {
		return err
	}

	return nil
}

// Update 更新环境
func (*env) Update(e *model.Env) error {
	return dao.Env.Update(e)
}

// Delete 删除环境
func (*env) Delete(id uint) error {
	_, has, err := dao.Service.Get(id)
	if err != nil {
		return err
	}

	if has {
		logger.Error("当前环境关联URL信息，请先删除关联信息")
		return errors.New("当前环境关联URL信息，请先删除关联信息")
	}

	return dao.Env.Delete(id)
}
