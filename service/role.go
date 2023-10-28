package service

import (
	"errors"
	"kubea-demo/dao"
	"kubea-demo/model"
)

var Roles r

type r struct{}

// List 返回用户列表
func (*r) List(roleName string, page, limit int) (*dao.Roles, error) {
	return dao.Role.List(roleName, page, limit)
}

// GetAll 查询所有页面信息
func (*r) GetAll() ([]*model.Role, error) {
	return dao.Role.GetAll()
}

// Add 新增
func (*r) Add(u *model.Role) error {
	_, has, err := dao.Role.Has(u.RoleName)
	if err != nil {
		return err
	}
	if has {
		return errors.New("该数据已存在，请重新添加")
	}

	//不存在则创建
	return dao.Role.Add(u)
}

// Update 更新
func (*r) Update(u *model.Role) error {
	return dao.Role.Update(u)
}

// Delete 删除
func (*r) Delete(id uint) error {
	return dao.Role.Delete(id)
}
