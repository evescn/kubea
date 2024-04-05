package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kubea/model"
	"kubea/service"
	"net/http"
)

var Env env

type env struct{}

// List 返回环境列表
func (*env) List(c *gin.Context) {
	params := new(struct {
		EnvName string `form:"env_name"`
		Page    int    `form:"page"`
		Limit   int    `form:"limit"`
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

	data, err := service.Env.List(params.EnvName, params.Page, params.Limit)
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
		"msg":  "获取环境列表成功",
		"data": data,
	})
}

// Add 创建环境
func (*env) Add(c *gin.Context) {
	params := new(model.Env)

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

	err := service.Env.Add(params)
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
		"msg":  "新增环境成功",
		"data": nil,
	})
}

// Update 更新环境
func (*env) Update(c *gin.Context) {
	params := new(model.Env)

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

	err := service.Env.Update(params)
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
		"msg":  "更新环境信息成功",
		"data": nil,
	})
}

// Delete 删除环境
func (*env) Delete(c *gin.Context) {
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

	err := service.Env.Delete(params.ID)
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
		"msg":  "删除环境成功",
		"data": nil,
	})
}
