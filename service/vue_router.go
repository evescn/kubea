package service

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"kubea-demo/dao"
)

var VueRouter vueRouter

type vueRouter struct{}

type Permission struct {
	Path     string        `json:"path"`
	Name     string        `json:"name"`
	Icon     string        `json:"icon"`
	Children []interface{} `json:"children"`
}

type Child struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Meta Meta   `json:"meta"`
}

type Meta struct {
	Title       string `json:"title"`
	RequireAuth bool   `json:"requireAuth"`
}

func (*vueRouter) SetRouter(roleID uint) ([]*Permission, error) {
	// 获取权限对应关系
	//logger.Info("roleID: ", roleID)
	// 超级管理员
	if roleID == 1 {
		return VueRouter.AdminRouter()
	}

	// 普通用户
	return VueRouter.UserRouter(roleID)
}

// AdminRouter 普通用户权限获取
func (*vueRouter) AdminRouter() ([]*Permission, error) {
	// 获取权限对应关系
	data, err := Menu.GetAll()
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		fmt.Println(item)
	}

	// 定义数据
	var (
		// 返回数据信息
		permissionList = make([]*Permission, 0)
		// 父节点信息
		permissionData = new(Permission)
		// 子父节点信息
		subPermissionData = new(Permission)
	)

	// 遍历数据，获数据信息
	for _, menuData := range data {
		permissionData = &Permission{
			Path:     menuData.Path,
			Name:     menuData.Name,
			Icon:     menuData.Icon,
			Children: nil,
		}
		//logger.Info("menuData: ", menuData)

		// 遍历1级页面
		subData, _, err := dao.SubMenu.GetP(menuData.ID)
		//logger.Info("subData: ", subData)
		if err != nil {
			return nil, err
		}
		if len(subData) == 0 {
			// 没有2级页面
			childData := &Child{
				Path: menuData.Path,
				Name: menuData.Name,
				Icon: menuData.Icon,
				Meta: Meta{
					Title:       menuData.Name,
					RequireAuth: true,
				},
			}
			permissionData.Children = append(permissionData.Children, childData)

			// 插入数据
			permissionList = append(permissionList, permissionData)
		} else {
			// 插入数据
			permissionList = append(permissionList, permissionData)

			// 遍历2级页面
			for _, subMenuData := range subData {
				//logger.Info("------------------------")
				//logger.Info("permissionData: ", permissionData)
				//logger.Info("subData: ", subData)
				subSubData, _, err := dao.SubSubMenu.GetP(subMenuData.ID)
				if err != nil {
					return nil, err
				}
				if len(subSubData) == 0 {
					// 没有2级页面
					childData := &Child{
						Path: subMenuData.Path,
						Name: subMenuData.Name,
						Icon: subMenuData.Icon,
						Meta: Meta{
							Title:       subMenuData.Name,
							RequireAuth: true,
						},
					}
					permissionData.Children = append(permissionData.Children, childData)

					// 插入数据
					//permissionList = append(permissionList, permissionData)
					//logger.Info("+++++++++++++++++++++++")
				} else {
					subPermissionData = &Permission{
						Path:     subMenuData.Path,
						Name:     subMenuData.Name,
						Icon:     subMenuData.Icon,
						Children: nil,
					}
					permissionData.Children = append(permissionData.Children, subPermissionData)

					// 遍历3级页面
					for _, subSubMenuData := range subSubData {
						fmt.Println("subSubData: ", subSubData)
						// 没有4级页面
						childData := &Child{
							Path: subSubMenuData.Path,
							Name: subSubMenuData.Name,
							Meta: Meta{
								Title:       subSubMenuData.Name,
								RequireAuth: true,
							},
						}
						subPermissionData.Children = append(subPermissionData.Children, childData)

					}
				}
			}
		}

	}

	return permissionList, nil
}

// UserRouter 普通用户权限获取
func (*vueRouter) UserRouter(roleID uint) ([]*Permission, error) {
	// 获取权限对应关系
	//logger.Info("roleID: ", roleID)
	data, err := RoleMenuRelation.Get(roleID)
	if err != nil {
		return nil, err
	}

	// 定义数据
	var (
		// 返回数据信息
		permissionList = make([]*Permission, 0)
		// 父节点信息
		permissionData = new(Permission)
		// 子父节点信息
		subPermissionData = new(Permission)
		// 判断是否已创建对应父节点
		menuList    = make(map[uint]bool)
		subMenuList = make(map[uint]bool)
	)

	// 遍历数据，获数据信息
	for _, item := range data {
		// 1级页面
		if _, ok := menuList[item.PageID]; !ok {
			menuData, has, err := Menu.Get(item.PageID)
			if err != nil {
				return nil, err
			}
			if !has {
				logger.Error("当前Menu数据不存在")
				return nil, errors.New("当前Menu数据不存在")
			}

			permissionData = &Permission{
				Path:     menuData.Path,
				Name:     menuData.Name,
				Icon:     menuData.Icon,
				Children: nil,
			}

			// 没有2级页面
			if item.SubPageID == 0 {
				childData := &Child{
					Path: menuData.Path,
					Name: menuData.Name,
					Icon: menuData.Icon,
					Meta: Meta{
						Title:       menuData.Name,
						RequireAuth: true,
					},
				}

				permissionData.Children = append(permissionData.Children, childData)
			}

			// 插入数据
			permissionList = append(permissionList, permissionData)
			menuList[item.PageID] = true

			// 没有2级页面，跳出本次循环，进行下一次循环
			if item.SubPageID == 0 {
				continue
			}
		}

		// 2级页面
		if _, ok := subMenuList[item.SubPageID]; !ok {
			// 获取2级菜单对应关系
			subMenuData, has, err := SubMenu.Get(item.SubPageID)
			if err != nil {
				return nil, err
			}
			if !has {
				logger.Error("当前SubMenu数据不存在")
				return nil, errors.New("当前SubMenu数据不存在")
			}

			// 没有3级页面
			if item.SubSubPageID == 0 {
				childData := &Child{
					Path: subMenuData.Path,
					Name: subMenuData.Name,
					Icon: subMenuData.Icon,
					Meta: Meta{
						Title:       subMenuData.Name,
						RequireAuth: true,
					},
				}
				permissionData.Children = append(permissionData.Children, childData)
			} else {
				subPermissionData = &Permission{
					Path:     subMenuData.Path,
					Name:     subMenuData.Name,
					Icon:     subMenuData.Icon,
					Children: nil,
				}

				permissionData.Children = append(permissionData.Children, subPermissionData)
			}

			subMenuList[item.SubPageID] = true

			// 没有3级页面，跳出本次循环，进行下一次循环
			if item.SubSubPageID == 0 {
				continue
			}

		}

		// 3级页面
		if item.SubSubPageID != 0 {
			// 获取2级菜单对应关系
			subSubMenuData, has, err := SubSubMenu.Get(item.SubSubPageID)
			if err != nil {
				return nil, err
			}
			if !has {
				logger.Error("当前SubMenu数据不存在")
				return nil, errors.New("当前SubMenu数据不存在")
			}

			// 没有4级页面
			childData := &Child{
				Path: subSubMenuData.Path,
				Name: subSubMenuData.Name,
				Meta: Meta{
					Title:       subSubMenuData.Name,
					RequireAuth: true,
				},
			}
			subPermissionData.Children = append(subPermissionData.Children, childData)
		}

	}

	return permissionList, nil
}
