package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubea-demo/model"
	"kubea-demo/service"
	"net/http"
)

var Deploy deploy

type deploy struct{}

// List 列表
func (*deploy) List(c *gin.Context) {
	//接收参数
	params := new(struct {
		En      string `form:"en"`
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
	data, err := service.Deploy.List(params.En, params.AppName, params.Page, params.Limit)
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
		"msg":  "获取部署列表成功",
		"data": data,
	})
}

// Update 更新
func (*deploy) Update(c *gin.Context) {
	//接收参数
	params := new(model.Deploy)

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
	err := service.Deploy.Update(params)
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
		"msg":  "更新部署信息成功",
		"data": nil,
	})
}

// Add 新增
func (*deploy) Add(c *gin.Context) {
	//接收参数
	params := new(model.Deploy)

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
	err := service.Deploy.Add(params)
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
		"msg":  "新增部署成功",
		"data": nil,
	})
}

// Delete 删除
func (*deploy) Delete(c *gin.Context) {
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
	err := service.Deploy.Delete(params.ID)
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
		"msg":  "删除部署信息成功",
		"data": nil,
	})
}

// GetLog 查看日志
func (*deploy) GetLog(c *gin.Context) {
	//接收参数
	params := new(struct {
		ID uint `form:"id"`
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
	data, _, err := service.Deploy.GetLog(params.ID)
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
		"msg":  "删除部署信息成功",
		"data": data,
	})
}
