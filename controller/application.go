package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubea-demo/model"
	"kubea-demo/service"
	"net/http"
)

var App app

type app struct{}

// List 列表
func (*app) List(c *gin.Context) {
	//接收参数
	params := new(struct {
		AppName string `form:"app_name"`
		Page    int    `form:"page"`
		Limit   int    `form:"limit"`
	})

	//绑定参数
	if err := c.Bind(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用Service方法
	data, err := service.App.List(params.AppName, params.Page, params.Limit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//返回
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取应用列表成功",
		"data": data,
	})
}

// GetAll 获取所有应用
func (*app) GetAll(c *gin.Context) {
	//调用Service方法
	data, err := service.App.GetAll()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//返回
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取所有应用信息成功",
		"data": data,
	})
}

// GetAllTags 查询所有tag
func (*app) GetAllTags(c *gin.Context) {
	//接收参数
	params := new(struct {
		AppId uint `form:"app_id"`
	})

	//绑定参数
	if err := c.Bind(params); err != nil {
		logger.Error("Bind请求参数失败," + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用Service方法
	data, err := service.App.GetAllTabByApp(params.AppId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//返回
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取应用Tag列表成功",
		"data": data,
	})
}

// Update 更新
func (*app) Update(c *gin.Context) {
	//接收参数
	params := new(model.App)

	//绑定参数
	if err := c.ShouldBind(params); err != nil {
		logger.Error("ShouldBind请求参数失败," + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用Service方法
	err := service.App.Update(params)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//返回
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "更新应用信息成功",
		"data": nil,
	})
}

// Add 新增
func (*app) Add(c *gin.Context) {
	//接收参数
	params := new(model.App)

	//绑定参数
	if err := c.ShouldBind(params); err != nil {
		logger.Error("ShouldBind请求参数失败," + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用Service方法
	err := service.App.Add(params)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//返回
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "新增应用成功",
		"data": nil,
	})
}

// Delete 删除
func (*app) Delete(c *gin.Context) {
	//接收参数
	params := new(struct {
		ID uint `json:"id"`
	})

	//绑定参数
	if err := c.ShouldBind(params); err != nil {
		logger.Error("ShouldBind请求参数失败," + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用Service方法
	err := service.App.Delete(params.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//返回
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除应用成功",
		"data": nil,
	})
}
