package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"kubea-demo/dao"
	"kubea-demo/model"
)

var User user

type user struct{}

// HashPassword 密码加密 在创建或更新用户时，对密码进行哈希
func (u *user) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword 用户登录验证密码
func (u *user) VerifyPassword(db, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(db), []byte(password))
	if err != nil {
		return false, errors.New("用户名或密码错误")
	}
	return true, nil
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

// Update 更新
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

// Delete 删除
func (*user) Delete(id uint) error {
	return dao.User.Delete(id)
}
