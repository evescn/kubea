package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"kubea-demo/db"
	"kubea-demo/model"
)

var User user

type user struct{}

// Has 根据用户名查询，用于代码层去重，查询账号信息
func (*user) Has(username string) (*model.User, bool, error) {
	data := new(model.User)
	tx := db.GORM.Where("username = ?", username).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据用户名查询User失败," + tx.Error.Error())
		return nil, false, errors.New("根据用户名查询User失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Update 更新
func (*user) Update(u *model.User) error {
	tx := db.GORM.Model(&model.User{}).Where("username = ?", u.UserName).Updates(&u)
	if tx.Error != nil {
		logger.Error("更新User信息失败," + tx.Error.Error())
		return errors.New("更新User信息失败," + tx.Error.Error())
	}

	return nil
}

// Add 新增
func (*user) Add(u *model.User) error {
	tx := db.GORM.Create(&u)
	if tx.Error != nil {
		logger.Error("新增User信息失败," + tx.Error.Error())
		return errors.New("新增User信息失败," + tx.Error.Error())
	}

	return nil
}

// Delete 删除
func (*user) Delete(id uint) error {
	data := new(model.User)
	data.ID = id
	tx := db.GORM.Delete(&data)
	if tx.Error != nil {
		logger.Error("删除User信息失败," + tx.Error.Error())
		return errors.New("删除User信息失败," + tx.Error.Error())
	}

	return nil
}
