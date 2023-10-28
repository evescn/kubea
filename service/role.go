package service

import (
	"errors"
	"kubea/dao"
	"kubea/model"
)

var Roles r

type r struct{}

// List 返回角色列表
func (*r) List(roleName string, page, limit int) (*dao.Roles, error) {
	return dao.Role.List(roleName, page, limit)
}

// GetAll 查询所有角色信息，添加账号选择角色
func (*r) GetAll() ([]*model.Role, error) {
	return dao.Role.GetAll()
}

// Add 新增角色
func (*r) Add(u *model.Role) error {
	_, has, err := dao.Role.Has(u.RoleName)
	if err != nil {
		return err
	}
	if has {
		return errors.New("该角色数据已存在，请重新添加")
	}

	//不存在则创建
	return dao.Role.Add(u)
}

// Update 更新角色
func (*r) Update(u *model.Role) error {
	return dao.Role.Update(u)
}

// Delete 删除角色
func (*r) Delete(id uint) error {
	return dao.Role.Delete(id)
}
