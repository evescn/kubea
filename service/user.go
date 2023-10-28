package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"kubea/dao"
	"kubea/model"
)

var User user

type user struct{}

// HashPassword 密码加密 在创建或更新用户时，对密码进行哈希
func (*user) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword 用户登录验证密码
func (*user) VerifyPassword(db, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(db), []byte(password))
	if err != nil {
		return false, errors.New("用户名或密码错误")
	}
	return true, nil
}

// List 返回用户列表
func (*user) List(userName string, page, limit int) (*dao.Users, error) {
	return dao.User.List(userName, page, limit)
}

// Add 新增
func (*user) Add(u *model.User) error {
	if u.UserName == "" || u.Password == "" {
		return errors.New("请填写用户名和密码信息")
	}
	_, has, err := dao.User.Has(u.UserName)
	if err != nil {
		return err
	}
	if has {
		return errors.New("该数据已存在，请重新添加")
	}

	//不存在则创建
	u.Password, err = User.HashPassword(u.Password)
	if err != nil {
		return err
	}

	return dao.User.Add(u)
}

// Update 更新密码
func (*user) Update(username, oldPassword, newPassword string) error {

	data, has, err := dao.User.Has(username)

	if err != nil {
		return err
	}
	if !has {
		return errors.New("更新密码失败, 用户名或密码错误")
	}

	_, err = User.VerifyPassword(data.Password, oldPassword)
	if err != nil {
		return err
	}

	u := new(model.User)
	u.UserName = username
	u.Password, err = User.HashPassword(newPassword)

	return dao.User.Update(u)
}

// UpdateAdmin 更新密码
func (*user) UpdateAdmin(username, newPassword string) error {
	u := new(model.User)
	u.UserName = username
	hashPassword, err := User.HashPassword(newPassword)
	if err != nil {
		return err
	}
	u.Password = hashPassword

	return dao.User.Update(u)
}

// Delete 删除
func (*user) Delete(id uint) error {
	return dao.User.Delete(id)
}

// UpdateRole 更新用户权限
func (*user) UpdateRole(username string, role uint) error {

	_, has, err := dao.User.Has(username)

	if err != nil {
		return err
	}
	if !has {
		return errors.New("更新权限失败, 用户信息错误")
	}

	u := new(model.User)
	u.UserName = username
	u.Role = role

	return dao.User.Update(u)
}
