package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubea-demo/model"
	"kubea-demo/service"
	"net/http"
)

var Roles role

type role struct{}

// List 返回环境列表
func (*role) List(c *gin.Context) {
	params := new(struct {
		RoleName string `form:"role_name"`
		Page     int    `form:"page"`
		Limit    int    `form:"limit"`
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

	data, err := service.Roles.List(params.RoleName, params.Page, params.Limit)
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
		"msg":  "获取角色列表成功",
		"data": data,
	})
}

// GetAll 所有角色信息，列表展示，新增用户获取角色数据
func (*role) GetAll(c *gin.Context) {

	//调用Service方法
	data, err := service.Roles.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//返回
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取所有角色信息成功",
		"data": data,
	})
}

// Add 新增
func (*role) Add(c *gin.Context) {
	//接收参数
	params := new(model.Role)

	//绑定参数
	if err := c.ShouldBind(params); err != nil {
		logger.Error("ShouldBind请求参数失败," + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用Service方法
	err := service.Roles.Add(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//返回
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "新增角色信息成功",
		"data": nil,
	})
}

// Update 更新
func (*role) Update(c *gin.Context) {
	//接收参数
	params := new(model.Role)

	//绑定参数
	if err := c.ShouldBind(params); err != nil {
		logger.Error("ShouldBind请求参数失败," + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用Service方法
	err := service.Roles.Update(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//返回
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "更新角色信息成功",
		"data": nil,
	})
}

// Delete 删除
func (*role) Delete(c *gin.Context) {
	//接收参数
	params := new(struct {
		ID uint `json:"id"`
	})

	//绑定参数
	if err := c.ShouldBind(params); err != nil {
		logger.Error("ShouldBind请求参数失败," + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用Service方法
	err := service.Roles.Delete(params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//返回
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除角色信息成功",
		"data": nil,
	})
}
