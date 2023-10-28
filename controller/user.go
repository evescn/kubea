package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubea-demo/model"
	"kubea-demo/service"
	"net/http"
)

var User user

type user struct{}

// List 返回环境列表
func (*user) List(c *gin.Context) {
	params := new(struct {
		UserName string `form:"username"`
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

	data, err := service.User.List(params.UserName, params.Page, params.Limit)
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
		"msg":  "获取用户列表成功",
		"data": data,
	})
}

// Add 新增
func (*user) Add(c *gin.Context) {
	//接收参数
	params := new(model.User)

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
	err := service.User.Add(params)
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
		"msg":  "新增用户成功",
		"data": nil,
	})
}

// Update 更新
func (*user) Update(c *gin.Context) {
	//接收参数
	params := new(struct {
		ID          uint   `json:"id"`
		UserName    string `json:"username"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
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
	err := service.User.Update(params.UserName, params.OldPassword, params.NewPassword)
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
		"msg":  "更新用户信息成功",
		"data": nil,
	})
}

// UpdateAdmin 更新
func (*user) UpdateAdmin(c *gin.Context) {
	//接收参数
	params := new(struct {
		ID          uint   `json:"id"`
		UserName    string `json:"username"`
		NewPassword string `json:"newPassword"`
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
	err := service.User.UpdateAdmin(params.UserName, params.NewPassword)
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
		"msg":  "更新用户信息成功",
		"data": nil,
	})
}

// Delete 删除
func (*user) Delete(c *gin.Context) {
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
	err := service.User.Delete(params.ID)
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
		"msg":  "删除用户成功",
		"data": nil,
	})
}

// UpdateRole 更新
func (*user) UpdateRole(c *gin.Context) {
	//接收参数
	params := new(struct {
		UserName string `json:"username"`
		Role     uint   `json:"role"`
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
	err := service.User.UpdateRole(params.UserName, params.Role)
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
		"msg":  "更新用户权限成功",
		"data": nil,
	})
}
