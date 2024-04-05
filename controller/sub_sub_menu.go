package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kubea/model"
	"kubea/service"
	"net/http"
)

var SubSubMenus subSubMenu

type subSubMenu struct{}

// List 返回3级菜单列表
func (*subSubMenu) List(c *gin.Context) {
	params := new(struct {
		SubSubMenuName string `form:"sub_sub_menu_name"`
		Page           int    `form:"page"`
		Limit          int    `form:"limit"`
	})

	// 绑定请求参数
	//绑定参数
	if err := c.Bind(params); err != nil {
		zap.L().Error("Bind 请求参数失败：" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.SubSubMenu.List(params.SubSubMenuName, params.Page, params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 90200,
		"msg":  "获取3级菜单列表成功",
		"data": data,
	})
}

// Add 创建3级菜单
func (*subSubMenu) Add(c *gin.Context) {
	params := new(model.SubSubMenu)

	// 绑定请求参数
	if err := c.ShouldBind(params); err != nil {
		zap.L().Error("Bind 请求参数失败：" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.SubSubMenu.Add(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 90200,
		"msg":  "新增3级菜单成功！",
		"data": nil,
	})
}

// Update 更新3级菜单
func (*subSubMenu) Update(c *gin.Context) {
	params := new(model.SubSubMenu)

	// 绑定请求参数
	if err := c.ShouldBind(params); err != nil {
		zap.L().Error("Bind 请求参数失败：" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.SubSubMenu.Update(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 90200,
		"msg":  "更新3级菜单信息成功！",
		"data": nil,
	})
}

// Delete 删除3级菜单
func (*subSubMenu) Delete(c *gin.Context) {
	params := new(struct {
		ID uint `json:"id"`
	})

	// 绑定请求参数
	if err := c.ShouldBind(params); err != nil {
		zap.L().Error("Bind 请求参数失败：" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.SubSubMenu.Delete(params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 90200,
		"msg":  "删除3级菜单成功！",
		"data": nil,
	})
}
