package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kubea/model"
	"kubea/service"
	"net/http"
)

var Menu menu

type menu struct{}

// List 返回1级菜单列表
func (*menu) List(c *gin.Context) {
	params := new(struct {
		MenuName string `form:"menu_name"`
		Page     int    `form:"page"`
		Limit    int    `form:"limit"`
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

	data, err := service.Menu.List(params.MenuName, params.Page, params.Limit)
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
		"msg":  "获取1级菜单列表成功",
		"data": data,
	})
}

// Add 创建1级菜单
func (*menu) Add(c *gin.Context) {
	params := new(model.Menu)

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

	err := service.Menu.Add(params)
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
		"msg":  "新增1级菜单成功",
		"data": nil,
	})
}

// Update 更新1级菜单
func (*menu) Update(c *gin.Context) {
	params := new(model.Menu)

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

	err := service.Menu.Update(params)
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
		"msg":  "更新1级菜单信息成功",
		"data": nil,
	})
}

// Delete 删除1级菜单
func (*menu) Delete(c *gin.Context) {
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

	err := service.Menu.Delete(params.ID)
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
		"msg":  "删除1级菜单成功",
		"data": nil,
	})
}
