package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"kubea-demo/db"
	"kubea-demo/model"
)

var RoleMenuRelation roleMenuRelation

type roleMenuRelation struct{}

//type RoleMenuRelations struct {
//	Items []*model.RoleMenuRelation `json:"items"`
//	Total int                       `json:"total"`
//}

// Get 根据 roleID 查询，查询账号权限信息
func (*roleMenuRelation) Get(roleID uint) ([]*model.RoleMenuRelation, error) {
	//定义返回值的内容
	roleMenuRelationList := make([]*model.RoleMenuRelation, 0)

	//数据库查询
	tx := db.GORM.Model(model.RoleMenuRelation{}).
		Where("role_id = ?", roleID).
		Order("page_id").
		Find(&roleMenuRelationList)
	if tx.Error != nil {
		logger.Error("根据RoleID查询RoleMenuRelation失败," + tx.Error.Error())
		return nil, errors.New("根据RoleID查询RoleMenuRelation失败," + tx.Error.Error())
	}

	return roleMenuRelationList, nil
}

// Has 根据用户名查询，用于代码层去重，查询账号信息
func (*roleMenuRelation) Has(id uint) (*model.RoleMenuRelation, bool, error) {
	data := new(model.RoleMenuRelation)
	tx := db.GORM.Where("id = ?", id).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		logger.Error("根据用户名查询RoleMenuRelation失败," + tx.Error.Error())
		return nil, false, errors.New("根据用户名查询RoleMenuRelation失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Add 新增
func (*roleMenuRelation) Add(u *model.RoleMenuRelation) error {
	tx := db.GORM.Create(&u)
	if tx.Error != nil {
		logger.Error("新增RoleMenuRelation信息失败," + tx.Error.Error())
		return errors.New("新增RoleMenuRelation信息失败," + tx.Error.Error())
	}

	return nil
}

// Update 更新
func (*roleMenuRelation) Update(u *model.RoleMenuRelation) error {
	tx := db.GORM.Model(&model.RoleMenuRelation{}).Where("id = ?", u.ID).Updates(&u)
	if tx.Error != nil {
		logger.Error("更新RoleMenuRelation信息失败," + tx.Error.Error())
		return errors.New("更新RoleMenuRelation信息失败," + tx.Error.Error())
	}

	return nil
}

// Delete 删除
func (*roleMenuRelation) Delete(u *model.RoleMenuRelation) error {
	tx := db.GORM.Where("role_id = ? and page_id =? and sub_page_id =? and sub_sub_page_id =? ", u.RoleID, u.PageID, u.SubPageID, u.SubSubPageID).Delete(&u)
	if tx.Error != nil {
		logger.Error("删除RoleMenuRelation信息失败," + tx.Error.Error())
		return errors.New("删除RoleMenuRelation信息失败," + tx.Error.Error())
	}

	return nil
}
