package service

import (
	"fmt"
	"kubea-demo/dao"
	"kubea-demo/model"
	"strconv"
	"strings"
)

var RoleMenuRelation roleMenuRelation

type roleMenuRelation struct{}

type TreeData struct {
	Name     string        `json:"name"`
	Disabled bool          `json:"disabled" `
	Children []interface{} `json:"children"`
}

type NodeData struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// GetAll 查询账号权限信息
func (*roleMenuRelation) GetAll() ([]*TreeData, error) {
	// 新建数据
	var tmpData = make([]*TreeData, 0)

	treeData := &TreeData{
		Name:     "菜单权限管理",
		Disabled: true,
		Children: nil,
	}

	//获取所有页面
	data, err := Menu.GetAll()
	if err != nil {
		return nil, err
	}
	for _, item := range data {
		// 查询2级页面确定当前节点是否为1级节点
		subData, _, err := dao.SubMenu.GetP(item.ID)
		if err != nil {
			return nil, err
		}

		if len(subData) == 0 {
			// 创建子节点
			nodeData := &NodeData{
				Name: item.Name,
				Key:  fmt.Sprintf("%v-%v-%v", item.ID, 0, 0),
			}
			treeData.Children = append(treeData.Children, nodeData)
		} else {
			subTreeData := &TreeData{
				Name:     item.Name,
				Disabled: true,
				Children: nil,
			}

			for _, subItem := range subData {
				// 查询3级页面确定当前节点是否为2级节点
				subSubData, _, err := dao.SubSubMenu.GetP(subItem.ID)
				if err != nil {
					return nil, err
				}

				if len(subSubData) == 0 {
					// 创建子节点
					nodeSubData := &NodeData{
						Name: subItem.Name,
						Key:  fmt.Sprintf("%v-%v-%v", item.ID, subItem.ID, 0),
					}
					subTreeData.Children = append(subTreeData.Children, nodeSubData)
				} else {
					subSubTreeData := &TreeData{
						Name:     subItem.Name,
						Disabled: true,
						Children: nil,
					}

					for _, subSubItem := range subSubData {
						// 创建子节点
						nodeSubSubData := &NodeData{
							Name: subSubItem.Name,
							Key:  fmt.Sprintf("%v-%v-%v", item.ID, subItem.ID, subSubItem.ID),
						}
						subSubTreeData.Children = append(subSubTreeData.Children, nodeSubSubData)
					}
					subTreeData.Children = append(subTreeData.Children, subSubTreeData)
				}

			}

			treeData.Children = append(treeData.Children, subTreeData)
		}

	}

	tmpData = append(tmpData, treeData)

	return tmpData, nil
}

// Get 根据 roleID 查询，查询账号权限信息
func (*roleMenuRelation) Get(roleID uint) ([]*model.RoleMenuRelation, error) {
	return dao.RoleMenuRelation.Get(roleID)
}

// GetPermissions 根据 roleID 查询，查询账号权限信息
func (*roleMenuRelation) GetPermissions(roleID uint) ([]string, error) {
	var rolePermissions = make([]string, 0)
	data, err := dao.RoleMenuRelation.Get(roleID)

	if err != nil {
		return nil, err
	}
	for _, item := range data {
		rolePermissions = append(rolePermissions, fmt.Sprintf("%v-%v-%v", item.PageID, item.SubPageID, item.SubSubPageID))
	}

	return rolePermissions, nil
}

// Update 更新环境
func (*roleMenuRelation) Update(roleID uint, newRoleMenuRelation, oldRoleMenuRelation []string) error {
	// 需要新增的数据
	addData := RoleMenuRelation.SliceDifference(oldRoleMenuRelation, newRoleMenuRelation)
	for _, item := range addData {
		arr := strings.Split(item, "-")
		pageID, _ := strconv.Atoi(arr[0])
		subPageID, _ := strconv.Atoi(arr[1])
		subSubPageID, _ := strconv.Atoi(arr[2])
		data := &model.RoleMenuRelation{
			RoleID:       roleID,
			PageID:       uint(pageID),
			SubPageID:    uint(subPageID),
			SubSubPageID: uint(subSubPageID),
		}

		err := dao.RoleMenuRelation.Add(data)
		if err != nil {
			return err
		}
	}
	// 需要删除的数据
	delData := RoleMenuRelation.SliceDifference(newRoleMenuRelation, oldRoleMenuRelation)
	for _, item := range delData {
		arr := strings.Split(item, "-")
		pageID, _ := strconv.Atoi(arr[0])
		subPageID, _ := strconv.Atoi(arr[1])
		subSubPageID, _ := strconv.Atoi(arr[2])
		data := &model.RoleMenuRelation{
			RoleID:       roleID,
			PageID:       uint(pageID),
			SubPageID:    uint(subPageID),
			SubSubPageID: uint(subSubPageID),
		}
		//logger.Info(data)
		err := dao.RoleMenuRelation.Delete(data)
		if err != nil {
			return err
		}
	}
	//logger.Info("addData:", addData)
	//logger.Info("delData:", delData)
	//return dao.RoleMenuRelation.Update(e)
	return nil
}

// SliceDifference 计算切片的差集
func (*roleMenuRelation) SliceDifference(slice1, slice2 []string) []string {
	// 将 slice1 转换为 map，以便进行快速查找
	sliceMap := make(map[string]struct{})
	for _, v := range slice1 {
		sliceMap[v] = struct{}{}
	}

	// 计算 slice2 中存在但在 slice1 中不存在的元素
	difference := make([]string, 0)
	for _, v := range slice2 {
		if _, exists := sliceMap[v]; !exists {
			difference = append(difference, v)
		}
	}

	//logger.Info("差集:", difference)
	return difference
}
