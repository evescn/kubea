package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubea-demo/model"
	"kubea-demo/service"
	"net/http"
)

var Build build

type build struct{}

// List 列表
func (*build) List(c *gin.Context) {
	//接收参数
	params := new(struct {
		RepoName string `form:"repo_name"`
		Page     int    `form:"page"`
		Limit    int    `form:"limit"`
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
	data, err := service.Build.List(params.RepoName, params.Page, params.Limit)
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
		"msg":  "获取编译列表成功",
		"data": data,
	})
}

// Add 新增
func (*build) Add(c *gin.Context) {
	//接收参数
	params := new(model.Build)

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
	err := service.Build.Add(params)
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
		"msg":  "新增编排信息成功",
		"data": nil,
	})
}

// Update 更新
func (*build) Update(c *gin.Context) {
	//接收参数
	params := new(model.Build)

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
	err := service.Build.Update(params)
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
		"msg":  "更新编排信息成功",
		"data": nil,
	})
}
