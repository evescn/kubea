package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubea/model"
	"kubea/service"
	"net/http"
)

var SubMenu subMenu

type subMenu struct{}

// List 返回2级菜单列表
func (*subMenu) List(c *gin.Context) {
	params := new(struct {
		SubMenuName string `form:"sub_menu_name"`
		Page        int    `form:"page"`
		Limit       int    `form:"limit"`
	})

	// 绑定请求参数
	//绑定参数
	if err := c.Bind(params); err != nil {
		logger.Error("Bind 请求参数失败：" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.SubMenu.List(params.SubMenuName, params.Page, params.Limit)
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
		"msg":  "获取2级菜单列表成功",
		"data": data,
	})
}

// Add 创建2级菜单
func (*subMenu) Add(c *gin.Context) {
	params := new(model.SubMenu)

	// 绑定请求参数
	if err := c.ShouldBind(params); err != nil {
		logger.Error("Bind 请求参数失败：" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.SubMenu.Add(params)
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
		"msg":  "新增2级菜单成功！",
		"data": nil,
	})
}

// Update 更新2级菜单
func (*subMenu) Update(c *gin.Context) {
	params := new(model.SubMenu)

	// 绑定请求参数
	if err := c.ShouldBind(params); err != nil {
		logger.Error("Bind 请求参数失败：" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.SubMenu.Update(params)
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
		"msg":  "更新2级菜单信息成功！",
		"data": nil,
	})
}

// Delete 删除2级菜单
func (*subMenu) Delete(c *gin.Context) {
	params := new(struct {
		ID uint `json:"id"`
	})

	// 绑定请求参数
	if err := c.ShouldBind(params); err != nil {
		logger.Error("Bind 请求参数失败：" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.SubMenu.Delete(params.ID)
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
		"msg":  "删除2级菜单成功！",
		"data": nil,
	})
}
