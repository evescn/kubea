package service

import (
	"errors"
	"go.uber.org/zap"
	"kubea/dao"
	"kubea/model"
)

var Menu menus

type menus struct{}

// List 返回1级菜单
func (*menus) List(menuName string, page, limit int) (*dao.Menus, error) {
	return dao.Menu.List(menuName, page, limit)
}

// GetAll 查询所有1级菜单信息
func (*menus) GetAll() ([]*model.Menu, error) {
	return dao.Menu.GetAll()
}

// Get 根据 ID 查询1级菜单信息
func (*menus) Get(ID uint) (*model.Menu, bool, error) {
	return dao.Menu.Get(ID)
}

// Add 创建1级菜单
func (*menus) Add(m *model.Menu) error {
	// 判断1级菜单是否存在
	_, has, err := dao.Menu.Has(m.Path)
	if err != nil {
		return err
	}

	if has {
		zap.L().
			Error("当前1级菜单数据已存在，请重新创建")
		return errors.New("当前1级菜单数据已存在，请重新创建")
	}

	// 不存在则创建数据
	if err = dao.Menu.Add(m); err != nil {
		return err
	}

	return nil
}

// Update 更新1级菜单
func (*menus) Update(m *model.Menu) error {
	return dao.Menu.Update(m)
}

// Delete 删除1级菜单
func (*menus) Delete(ID uint) error {
	_, has, err := dao.SubMenu.GetP(ID)
	if err != nil {
		return err
	}

	if has {
		zap.L().Error("当前1级菜单页面关联子页面信息，请先删除关联信息")
		return errors.New("当前1级菜单页面关联子页面信息，请先删除关联信息")
	}

	return dao.Menu.Delete(ID)
}
