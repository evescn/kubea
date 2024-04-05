package service

import (
	"errors"
	"go.uber.org/zap"
	"kubea/dao"
	"kubea/model"
)

var SubSubMenu subSubMenu

type subSubMenu struct{}

// List 返回环境列表
func (*subSubMenu) List(subSubMenuName string, page, limit int) (*dao.SubSubMenus, error) {
	return dao.SubSubMenu.List(subSubMenuName, page, limit)
}

// Get 根据 ID 查询，查询2级菜单信息
func (*subSubMenu) Get(ID uint) (*model.SubSubMenu, bool, error) {
	return dao.SubSubMenu.Get(ID)
}

// Add 创建3级菜单
func (*subSubMenu) Add(m *model.SubSubMenu) error {
	// 判断3级菜单是否存在
	_, has, err := dao.SubSubMenu.Has(m.Path)
	if err != nil {
		return err
	}

	if has {
		zap.L().
			Error("当前3级菜单数据已存在，请重新创建")
		return errors.New("当前3级菜单数据已存在，请重新创建")
	}

	// 不存在则创建数据
	if err = dao.SubSubMenu.Add(m); err != nil {
		return err
	}

	return nil
}

// Update 更新3级菜单
func (*subSubMenu) Update(m *model.SubSubMenu) error {
	return dao.SubSubMenu.Update(m)
}

// Delete 删除3级菜单
func (*subSubMenu) Delete(id uint) error {
	return dao.SubSubMenu.Delete(id)
}
