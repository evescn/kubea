package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"kubea-demo/db"
	"kubea-demo/model"
)

var SubMenu subMenu

type subMenu struct{}

type SubMenus struct {
	Items []*model.SubMenu `json:"items"`
	Total int              `json:"total"`
}

// List 列表
// subMenuName用于模糊查询，过滤
// page，limit用于分页
func (*subMenu) List(subMenuName string, page, limit int) (*SubMenus, error) {
	//计算分页
	startSet := (page - 1) * limit

	//定义返回值的内容
	var (
		subMenuList = make([]*model.SubMenu, 0)
		total       = 0
	)

	//数据库查询，先查total
	tx := db.GORM.Model(model.SubMenu{}).
		Where("name like ? ", "%"+subMenuName+"%").
		Count(&total)

	if tx.Error != nil {
		logger.Error("获取SubMenu列表失败," + tx.Error.Error())
		return nil, errors.New("获取SubMenu列表失败," + tx.Error.Error())
	}

	//数据库查询，再查数据
	//当limit=10， total一定是10，因为count会在过滤和分页后执行
	tx = db.GORM.Model(model.SubMenu{}).
		Where("name like ?", "%"+subMenuName+"%").
		Limit(limit).
		Offset(startSet).
		Order("name").
		Find(&subMenuList)
	if tx.Error != nil {
		logger.Error("获取SubMenu列表失败," + tx.Error.Error())
		return nil, errors.New("获取SubMenu列表失败," + tx.Error.Error())
	}

	return &SubMenus{
		Items: subMenuList,
		Total: total,
	}, nil
}

// Get 根据 ID 查询，查询2级菜单信息
func (*subMenu) Get(ID uint) (*model.SubMenu, bool, error) {
	data := new(model.SubMenu)
	tx := db.GORM.Where("id = ?", ID).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据ID查询SubMenu失败," + tx.Error.Error())
		return nil, false, errors.New("根据ID查询SubMenu失败," + tx.Error.Error())
	}

	return data, true, nil
}

// GetP 根据 ParentID 查询，查询父节点下面的子节点信息
func (*subMenu) GetP(id uint) ([]*model.SubMenu, bool, error) {
	data := make([]*model.SubMenu, 0)
	tx := db.GORM.Model(&model.SubMenu{}).Where("parent_page_id = ?", id).Find(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据 ParentID 查询SubMenu失败," + tx.Error.Error())
		return nil, false, errors.New("根据 ParentID 查询SubMenu失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Has 根据路径查询，用于代码层去重，查询2级菜单信息
func (*subMenu) Has(pagePath string) (*model.SubMenu, bool, error) {
	data := new(model.SubMenu)
	tx := db.GORM.Where("path = ?", pagePath).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据路径查询SubMenu失败," + tx.Error.Error())
		return nil, false, errors.New("根据路径查询SubMenu失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Update 更新
func (*subMenu) Update(u *model.SubMenu) error {
	tx := db.GORM.Model(&model.SubMenu{}).Where("path = ?", u.Path).Updates(&u)
	if tx.Error != nil {
		logger.Error("更新SubMenu信息失败," + tx.Error.Error())
		return errors.New("更新SubMenu信息失败," + tx.Error.Error())
	}

	return nil
}

// Add 新增
func (*subMenu) Add(u *model.SubMenu) error {
	tx := db.GORM.Create(&u)
	if tx.Error != nil {
		logger.Error("新增SubMenu信息失败," + tx.Error.Error())
		return errors.New("新增SubMenu信息失败," + tx.Error.Error())
	}

	return nil
}

// Delete 删除
func (*subMenu) Delete(id uint) error {
	data := new(model.SubMenu)
	data.ID = id
	tx := db.GORM.Delete(&data)
	if tx.Error != nil {
		logger.Error("删除SubMenu信息失败," + tx.Error.Error())
		return errors.New("删除SubMenu信息失败," + tx.Error.Error())
	}

	return nil
}
