package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubea-demo/dao"
	"kubea-demo/service"
	"net/http"
)

var Service svc

type svc struct{}

// List 返回环境列表
func (*svc) List(c *gin.Context) {
	params := new(struct {
		Username    string `form:"username"`
		Role        uint   `form:"role"`
		ServiceName string `form:"svc_name"`
		Eid         uint   `form:"eid"`
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

	data, err := service.Services.List(params.Username, params.ServiceName, params.Role, params.Eid, params.Page, params.Limit)
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
		"msg":  "获取URL列表成功",
		"data": data,
	})
}

// Add 创建环境
func (*svc) Add(c *gin.Context) {
	params := new(dao.ServiceRes)

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

	err := service.Services.Add(params)
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
		"msg":  "新增URL信息成功",
		"data": nil,
	})
}

// Update 更新环境
func (*svc) Update(c *gin.Context) {
	params := new(dao.ServiceRes)

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

	err := service.Services.Update(params)
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
		"msg":  "更新URL信息成功",
		"data": nil,
	})
}

// Delete 删除服务
func (*svc) Delete(c *gin.Context) {
	params := new(struct {
		SID uint `json:"id"`
		PID uint `json:"pid"`
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

	err := service.Services.Delete(params.SID, params.PID)
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
		"msg":  "删除URL成功",
		"data": nil,
	})
}
