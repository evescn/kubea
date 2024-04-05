package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"kubea/db"
	"kubea/model"
)

var Role role

type role struct{}

type Roles struct {
	Items []*RoleTemp `json:"items"`
	Total int         `json:"total"`
}

type RoleTemp struct {
	*model.Role
	Key uint `json:"key"`
}

// List 列表
// roleName用于模糊查询，过滤
// page，limit用于分页
func (*role) List(roleName string, page, limit int) (*Roles, error) {
	//计算分页
	startSet := (page - 1) * limit

	//定义返回值的内容
	var (
		roleListTemp = make([]*model.Role, 0)
		roleList     = make([]*RoleTemp, 0)
		total        = 0
	)

	// 过滤超级管理员
	query := db.GORM.Model(&model.Role{}).
		Where("role_name != ?", "超级管理员")
	if roleName != "" {
		query = query.Where("role_name like ? ", "%"+roleName+"%")
	}

	tx := query.Count(&total)
	if tx.Error != nil {
		zap.L().Error("获取Role列表失败," + tx.Error.Error())
		return nil, errors.New("获取Role列表失败," + tx.Error.Error())
	}

	//分页数据
	tx = query.Limit(limit).
		Offset(startSet).
		Order("id").
		Find(&roleListTemp)

	for _, item := range roleListTemp {
		tmp := &RoleTemp{
			Role: item,
			Key:  item.ID,
		}
		roleList = append(roleList, tmp)
	}

	return &Roles{
		Items: roleList,
		Total: total,
	}, nil
}

// GetAll 查询所有路由权限信息
func (*role) GetAll() ([]*model.Role, error) {
	data := make([]*model.Role, 0)
	tx := db.GORM.Find(&data)
	if tx.Error != nil {
		zap.L().Error("查询所有Role失败," + tx.Error.Error())
		return nil, errors.New("查询所有Role失败," + tx.Error.Error())
	}

	return data, nil
}

// Get 根据 ID 查询，查询账号信息
func (*role) Get(id uint) (*model.Role, bool, error) {
	data := new(model.Role)
	tx := db.GORM.Where("ID = ?", id).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		zap.L().Error("根据用户名查询Role失败," + tx.Error.Error())
		return nil, false, errors.New("根据用户名查询Role失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Has 根据用户名查询，用于代码层去重，查询账号信息
func (*role) Has(roleName string) (*model.Role, bool, error) {
	data := new(model.Role)
	tx := db.GORM.Where("role_name = ?", roleName).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		zap.L().Error("根据用户名查询Role失败," + tx.Error.Error())
		return nil, false, errors.New("根据用户名查询Role失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Update 更新
func (*role) Update(u *model.Role) error {
	tx := db.GORM.Model(&model.Role{}).Where("role_name = ?", u.RoleName).Updates(&u)
	if tx.Error != nil {
		zap.L().Error("更新Role信息失败," + tx.Error.Error())
		return errors.New("更新Role信息失败," + tx.Error.Error())
	}

	return nil
}

// Add 新增
func (*role) Add(u *model.Role) error {
	tx := db.GORM.Create(&u)
	if tx.Error != nil {
		zap.L().Error("新增Role信息失败," + tx.Error.Error())
		return errors.New("新增Role信息失败," + tx.Error.Error())
	}

	return nil
}

// Delete 删除
func (*role) Delete(id uint) error {
	data := new(model.Role)
	data.ID = id
	tx := db.GORM.Delete(&data)
	if tx.Error != nil {
		zap.L().Error("删除Role信息失败," + tx.Error.Error())
		return errors.New("删除Role信息失败," + tx.Error.Error())
	}

	return nil
}
