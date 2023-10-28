package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"kubea-demo/db"
	"kubea-demo/model"
)

var SubSubMenu subSubMenu

type subSubMenu struct{}

type SubSubMenus struct {
	Items []*model.SubSubMenu `json:"items"`
	Total int                 `json:"total"`
}

// List 列表
// subSubMenuName用于模糊查询，过滤
// page，limit用于分页
func (*subSubMenu) List(subSubMenuName string, page, limit int) (*SubSubMenus, error) {
	//计算分页
	startSet := (page - 1) * limit

	//定义返回值的内容
	var (
		subSubMenuList = make([]*model.SubSubMenu, 0)
		total          = 0
	)

	//数据库查询，先查total
	tx := db.GORM.Model(model.SubSubMenu{}).
		Where("name like ? ", "%"+subSubMenuName+"%").
		Count(&total)

	if tx.Error != nil {
		logger.Error("获取SubSubMenu列表失败," + tx.Error.Error())
		return nil, errors.New("获取SubSubMenu列表失败," + tx.Error.Error())
	}

	//数据库查询，再查数据
	//当limit=10， total一定是10，因为count会在过滤和分页后执行
	tx = db.GORM.Model(model.SubSubMenu{}).
		Where("name like ?", "%"+subSubMenuName+"%").
		Limit(limit).
		Offset(startSet).
		Order("name").
		Find(&subSubMenuList)
	if tx.Error != nil {
		logger.Error("获取SubSubMenu列表失败," + tx.Error.Error())
		return nil, errors.New("获取SubSubMenu列表失败," + tx.Error.Error())
	}

	return &SubSubMenus{
		Items: subSubMenuList,
		Total: total,
	}, nil
}

// Get 根据 roleID 查询，查询3级菜单信息
func (*subSubMenu) Get(ID uint) (*model.SubSubMenu, bool, error) {
	data := new(model.SubSubMenu)
	tx := db.GORM.Where("id = ?", ID).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据ID查询SubSubMenu失败," + tx.Error.Error())
		return nil, false, errors.New("根据ID查询SubSubMenu失败," + tx.Error.Error())
	}

	return data, true, nil
}

// GetP 根据 ParentSubID 查询，查询二级页面节点下面的子节点信息
func (*subSubMenu) GetP(id uint) ([]*model.SubSubMenu, bool, error) {
	data := make([]*model.SubSubMenu, 0)
	tx := db.GORM.Model(&model.SubSubMenu{}).Where("parent_sub_page_id = ?", id).Find(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据 ParentID 查询SubMenu失败," + tx.Error.Error())
		return nil, false, errors.New("根据 ParentID 查询SubMenu失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Has 根据路径查询，用于代码层去重，查询3级菜单信息
func (*subSubMenu) Has(pagePath string) (*model.SubSubMenu, bool, error) {
	data := new(model.SubSubMenu)
	tx := db.GORM.Where("path = ?", pagePath).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据路径查询SubSubMenu失败," + tx.Error.Error())
		return nil, false, errors.New("根据路径查询SubSubMenu失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Update 更新
func (*subSubMenu) Update(u *model.SubSubMenu) error {
	tx := db.GORM.Model(&model.SubSubMenu{}).Where("path = ?", u.Path).Updates(&u)
	if tx.Error != nil {
		logger.Error("更新SubSubMenu信息失败," + tx.Error.Error())
		return errors.New("更新SubSubMenu信息失败," + tx.Error.Error())
	}

	return nil
}

// Add 新增
func (*subSubMenu) Add(u *model.SubSubMenu) error {
	tx := db.GORM.Create(&u)
	if tx.Error != nil {
		logger.Error("新增SubSubMenu信息失败," + tx.Error.Error())
		return errors.New("新增SubSubMenu信息失败," + tx.Error.Error())
	}

	return nil
}

// Delete 删除
func (*subSubMenu) Delete(id uint) error {
	data := new(model.SubSubMenu)
	data.ID = id
	//tx := db.GORM.Unscoped().Delete(&data)
	tx := db.GORM.Delete(&data)
	if tx.Error != nil {
		logger.Error("删除SubSubMenu信息失败," + tx.Error.Error())
		return errors.New("删除SubSubMenu信息失败," + tx.Error.Error())
	}

	return nil
}
