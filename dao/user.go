package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"kubea/db"
	"kubea/model"
)

var User user

type user struct{}

type Users struct {
	Items []*model.User `json:"items"`
	Total int           `json:"total"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Role     uint   `json:"role"`
}

// List 列表
// userName用于模糊查询，过滤
// page，limit用于分页
func (*user) List(userName string, page, limit int) (*Users, error) {
	//计算分页
	startSet := (page - 1) * limit

	//定义返回值的内容
	var (
		userList     = make([]*model.User, 0)
		userListInfo = make([]*UserInfo, 0)
		total        = 0
	)

	//数据库查询，先查total
	tx := db.GORM.Model(model.User{}).
		Where("username like ? ", "%"+userName+"%").
		Count(&total)

	if tx.Error != nil {
		zap.L().Error("获取User列表失败," + tx.Error.Error())
		return nil, errors.New("获取User列表失败," + tx.Error.Error())
	}

	//数据库查询，再查数据
	//当limit=10， total一定是10，因为count会在过滤和分页后执行
	tx = db.GORM.Model(model.User{}).
		Where("username like ?", "%"+userName+"%").
		Limit(limit).
		Offset(startSet).
		Order("username").
		Find(&userList)
	if tx.Error != nil {
		zap.L().Error("获取User列表失败," + tx.Error.Error())
		return nil, errors.New("获取User列表失败," + tx.Error.Error())
	}

	// 返回第一条记录即可
	for _, item := range userList {
		tmp := &UserInfo{
			ID:       item.ID,
			UserName: item.UserName,
			Role:     item.Role,
		}
		userListInfo = append(userListInfo, tmp)
	}

	return &Users{
		Items: userList,
		Total: total,
	}, nil
}

// Has 根据用户名查询，用于代码层去重，查询账号信息
func (*user) Has(userName string) (*model.User, bool, error) {
	data := new(model.User)
	tx := db.GORM.Where("username = ?", userName).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		zap.L().Error("根据用户名查询User失败," + tx.Error.Error())
		return nil, false, errors.New("根据用户名查询User失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Update 更新
func (*user) Update(u *model.User) error {
	tx := db.GORM.Model(&model.User{}).Where("username = ?", u.UserName).Updates(&u)
	if tx.Error != nil {
		zap.L().Error("更新User信息失败," + tx.Error.Error())
		return errors.New("更新User信息失败," + tx.Error.Error())
	}

	return nil
}

// Add 新增
func (*user) Add(u *model.User) error {
	tx := db.GORM.Create(&u)
	if tx.Error != nil {
		zap.L().Error("新增User信息失败," + tx.Error.Error())
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
		zap.L().Error("删除User信息失败," + tx.Error.Error())
		return errors.New("删除User信息失败," + tx.Error.Error())
	}

	return nil
}
