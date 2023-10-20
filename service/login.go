package service

import (
	"errors"
	"github.com/wonderivan/logger"
	"kubea-demo/config"
	"kubea-demo/dao"
)

var Login login

type login struct{}

// Auth 验证账号密码
func (l *login) Auth(username, password string) (err error) {
	if username == config.AdminUser {
		if password != config.AdminPwd {
			logger.Error("登录失败, 用户名或密码错误")
			return errors.New("登录失败, 用户名或密码错误")
		}
	} else {

		data, has, err := dao.User.Has(username)

		if err != nil {
			return err
		}
		if !has {
			return errors.New("登录失败, 用户名或密码错误")
		}

		_, err = User.VerifyPassword(data.Password, password)
		if err != nil {
			return err
		}

	}
	//if username == config.AdminUser && password == config.AdminPwd {
	//	return nil
	//} else {
	//	logger.Error("登录失败, 用户名或密码错误")
	//	return errors.New("登录失败, 用户名或密码错误")
	//}
	return nil
}
