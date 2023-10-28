package service

import (
	"errors"
	"github.com/wonderivan/logger"
	"kubea-demo/dao"
	"kubea-demo/model"
)

var SubSubMenu subSubMenu

type subSubMenu struct{}

// List 返回环境列表
func (*subSubMenu) List(subSubMenuName string, page, limit int) (*dao.SubSubMenus, error) {
	return dao.SubSubMenu.List(subSubMenuName, page, limit)
}

// Get 根据 ID 查询，查询页面信息
func (*subSubMenu) Get(ID uint) (*model.SubSubMenu, bool, error) {
	return dao.SubSubMenu.Get(ID)
}

// Add 创建环境
func (*subSubMenu) Add(m *model.SubSubMenu) error {
	// 判断环境是否存在
	_, has, err := dao.SubSubMenu.Has(m.Path)
	if err != nil {
		return err
	}

	if has {
		logger.Error("当前数据已存在，请重新创建")
		return errors.New("当前数据已存在，请重新创建")
	}

	// 不存在则创建数据
	if err = dao.SubSubMenu.Add(m); err != nil {
		return err
	}

	return nil
}

// Update 更新环境
func (*subSubMenu) Update(m *model.SubSubMenu) error {
	return dao.SubSubMenu.Update(m)
}

// Delete 删除环境
func (*subSubMenu) Delete(id uint) error {
	return dao.SubSubMenu.Delete(id)
}
