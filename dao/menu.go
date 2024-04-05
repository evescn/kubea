package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"kubea/db"
	"kubea/model"
)

var Menu menu

type menu struct{}

type Menus struct {
	Items []*model.Menu `json:"items"`
	Total int           `json:"total"`
}

// List 列表
// menuName用于模糊查询，过滤
// page，limit用于分页
func (*menu) List(menuName string, page, limit int) (*Menus, error) {
	//计算分页
	startSet := (page - 1) * limit

	//定义返回值的内容
	var (
		menuList = make([]*model.Menu, 0)
		total    = 0
	)

	//数据库查询，先查total
	tx := db.GORM.Model(model.Menu{}).
		Where("name like ? ", "%"+menuName+"%").
		Count(&total)

	if tx.Error != nil {
		zap.L().Error("获取Menu列表失败," + tx.Error.Error())
		return nil, errors.New("获取Menu列表失败," + tx.Error.Error())
	}

	//数据库查询，再查数据
	//当limit=10， total一定是10，因为count会在过滤和分页后执行
	tx = db.GORM.Model(model.Menu{}).
		Where("name like ?", "%"+menuName+"%").
		Limit(limit).
		Offset(startSet).
		Order("name").
		Find(&menuList)
	if tx.Error != nil {
		zap.L().Error("获取Menu列表失败," + tx.Error.Error())
		return nil, errors.New("获取Menu列表失败," + tx.Error.Error())
	}

	return &Menus{
		Items: menuList,
		Total: total,
	}, nil
}

// GetAll 查询所有1级菜单信息
func (*menu) GetAll() ([]*model.Menu, error) {
	data := make([]*model.Menu, 0)
	tx := db.GORM.Find(&data)
	if tx.Error != nil {
		zap.L().Error("查询所有Menu失败," + tx.Error.Error())
		return nil, errors.New("查询所有Menu失败," + tx.Error.Error())
	}

	return data, nil
}

// Get 根据 ID 查询，查询1级菜单信息
func (*menu) Get(ID uint) (*model.Menu, bool, error) {
	data := new(model.Menu)
	tx := db.GORM.Where("id = ?", ID).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		zap.L().Error("根据路径查询Menu失败," + tx.Error.Error())
		return nil, false, errors.New("根据路径查询Menu失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Has 根据路径查询，用于代码层去重，查询1级菜单信息
func (*menu) Has(pagePath string) (*model.Menu, bool, error) {
	data := new(model.Menu)
	tx := db.GORM.Where("path = ?", pagePath).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		zap.L().Error("根据路径查询Menu失败," + tx.Error.Error())
		return nil, false, errors.New("根据路径查询Menu失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Add 新增
func (*menu) Add(u *model.Menu) error {
	tx := db.GORM.Create(&u)
	if tx.Error != nil {
		zap.L().Error("新增Menu信息失败," + tx.Error.Error())
		return errors.New("新增Menu信息失败," + tx.Error.Error())
	}

	return nil
}

// Update 更新
func (*menu) Update(u *model.Menu) error {
	tx := db.GORM.Model(&model.Menu{}).Where("path = ?", u.Path).Updates(&u)
	if tx.Error != nil {
		zap.L().Error("更新Menu信息失败," + tx.Error.Error())
		return errors.New("更新Menu信息失败," + tx.Error.Error())
	}

	return nil
}

// Delete 删除
func (*menu) Delete(id uint) error {
	data := new(model.Menu)
	data.ID = id
	tx := db.GORM.Delete(&data)
	if tx.Error != nil {
		zap.L().Error("删除Menu信息失败," + tx.Error.Error())
		return errors.New("删除Menu信息失败," + tx.Error.Error())
	}

	return nil
}
