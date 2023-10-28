package service

import (
	"errors"
	"github.com/wonderivan/logger"
	"kubea-demo/config"
	"kubea-demo/dao"
	"kubea-demo/model"
)

var Login login

type login struct{}

// Auth 验证账号密码
func (l *login) Auth(username, password string) (*model.User, error) {
	if username == config.AdminUser {
		if password != config.AdminPwd {
			logger.Error("登录失败, 用户名或密码错误")
			return nil, errors.New("登录失败, 用户名或密码错误")
		} else {
			return &model.User{
				Role: 1,
			}, nil
		}
	} else {

		data, has, err := dao.User.Has(username)

		if err != nil {
			return nil, err
		}
		if !has {
			return nil, errors.New("登录失败, 用户名或密码错误")
		}

		_, err = User.VerifyPassword(data.Password, password)
		if err != nil {
			return nil, err
		}
		return data, nil

	}
	return nil, nil
}
