package service

import (
	"errors"
	"github.com/wonderivan/logger"
	"kubea/dao"
	"kubea/model"
)

var Services service

type service struct{}

// List 返回环境列表
func (*service) List(userName, serviceName string, role, eid uint, page, limit int) (*dao.Services, error) {
	// 判断用户权限，确实是否返回账户名密码信息
	data, has, err := dao.User.Has(userName)
	if err != nil {
		return nil, err
	}

	// 超级管理员
	if !has {
		return dao.Service.List(serviceName, eid, page, limit)
	}

	if data.Role != role {
		logger.Error("当前用户信息不存在，")
		return nil, errors.New("当前用户信息不存在")
	}

	roleData, has, err := dao.Role.Get(role)
	if !has {
		logger.Error("当前角色信息不存在，")
		return nil, errors.New("当前角色信息不存在")
	}

	// 运维用户也返回账户名密码信息
	if roleData.RoleName == "运维用户" {
		return dao.Service.List(serviceName, eid, page, limit)
	}

	return dao.Service.ListNoUserInfo(serviceName, eid, page, limit)
}

// Add 创建环境
func (*service) Add(se *dao.ServiceRes) error {
	// 判断环境是否存在
	_, has, err := dao.Service.Has(se.Name, se.Eid)
	if err != nil {
		return err
	}

	if has {
		logger.Error("当前数据已存在，请重新创建")
		return errors.New("当前数据已存在，请重新创建")
	}

	// 不存在则创建数据
	// 创建 password 数据
	p := &model.Password{
		PName:    se.PName,
		Password: se.Password,
	}

	pid, err := dao.Password.Add(p)
	if err != nil {
		return err
	}

	// 创建 service 数据
	s := &model.Service{
		Name:        se.Name,
		Url:         se.Url,
		Description: se.Description,
		Eid:         se.Eid,
		Pid:         pid,
	}

	if err = dao.Service.Add(s); err != nil {
		return err
	}

	return nil
}

// Update 更新环境
func (*service) Update(se *dao.ServiceRes) error {
	// 创建 password 数据
	p := &model.Password{
		ID:       se.Pid,
		PName:    se.PName,
		Password: se.Password,
	}
	//logger.Info(*p)
	if err := dao.Password.Update(p); err != nil {
		return err
	}

	// 创建 service 数据
	s := &model.Service{
		ID:          se.ID,
		Name:        se.Name,
		Url:         se.Url,
		Description: se.Description,
	}
	//logger.Info(*s)
	if err := dao.Service.Update(s); err != nil {
		return err
	}

	return nil
}

// Delete 删除环境
func (*service) Delete(sid, pid uint) error {
	// 删除 服务账号信息
	if err := dao.Password.Delete(pid); err != nil {
		return err
	}

	return dao.Service.Delete(sid)

}
