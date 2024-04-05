package service

import (
	"errors"
	"go.uber.org/zap"
	"kubea/dao"
	"kubea/model"
)

var SubMenu subMenu

type subMenu struct{}

// List 返回2级菜单列表
func (*subMenu) List(subMenuName string, page, limit int) (*dao.SubMenus, error) {
	return dao.SubMenu.List(subMenuName, page, limit)
}

// Get 根据 ID 查询，查询页面信息
func (*subMenu) Get(ID uint) (*model.SubMenu, bool, error) {
	return dao.SubMenu.Get(ID)
}

// Add 创建2级菜单
func (*subMenu) Add(m *model.SubMenu) error {
	// 判断2级菜单是否存在
	_, has, err := dao.SubMenu.Has(m.Path)
	if err != nil {
		return err
	}

	if has {
		zap.L().Error("当前2级菜单数据已存在，请重新创建")
		return errors.New("当前2级菜单数据已存在，请重新创建")
	}

	// 不存在则创建数据
	if err = dao.SubMenu.Add(m); err != nil {
		return err
	}

	return nil
}

// Update 更新2级菜单
func (*subMenu) Update(m *model.SubMenu) error {
	return dao.SubMenu.Update(m)
}

// Delete 删除2级菜单
func (*subMenu) Delete(ID uint) error {
	_, has, err := dao.SubSubMenu.GetP(ID)
	if err != nil {
		return err
	}

	if has {
		zap.L().Error("当前2级菜单页面关联子页面信息，请先删除关联信息")
		return errors.New("当前2级菜单页面关联子页面信息，请先删除关联信息")
	}

	return dao.SubMenu.Delete(ID)
}
