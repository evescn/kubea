package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubea-demo/service"
	"net/http"
)

var Servicev1 servicev1

type servicev1 struct{}

// GetServices 获取 Service 列表
func (s *servicev1) GetServices(c *gin.Context) {
	//client *kubernetes.Clientset, filterName, namespace string, limit, page int
	//接收参数,匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
		Cluster    string `form:"cluster"`
	})

	//绑定参数
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := c.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败， %v\n", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败， %v\n", err),
			"data": nil,
		})
		return
	}

	//获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用service方法，获取列表
	data, err := service.Servicev1.GetServices(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Service列表成功",
		"data": data,
	})

}

// GetServiceDetail 获取 Service 详情
func (s *servicev1) GetServiceDetail(c *gin.Context) {
	//client *kubernetes.Clientset, ServiceName, namespace string
	//接收参数,匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		ServiceName string `form:"svc_name"`
		Namespace   string `form:"namespace"`
		Cluster     string `form:"cluster"`
	})

	//绑定参数
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := c.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败， %v\n", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败， %v\n", err),
			"data": nil,
		})
		return
	}

	//获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用service方法，获取列表
	data, err := service.Servicev1.GetServiceDetail(client, params.ServiceName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Service详情成功",
		"data": data,
	})

}

// DeleteService 删除 Service
func (s *servicev1) DeleteService(c *gin.Context) {
	//client *kubernetes.Clientset, ServiceName, namespace string
	//接收参数,匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		ServiceName string `json:"svc_name"`
		Namespace   string `json:"namespace"`
		Cluster     string `json:"cluster"`
	})

	//绑定参数
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败， %v\n", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败， %v\n", err),
			"data": nil,
		})
		return
	}

	//获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用service方法，获取列表
	err = service.Servicev1.DeleteService(client, params.ServiceName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除Service成功",
		"data": nil,
	})

}

// UpdateService 更新 Service
func (s *servicev1) UpdateService(c *gin.Context) {
	//client *kubernetes.Clientset, content, namespace string
	//接收参数,匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		Content   string `json:"content"`
		Namespace string `json:"namespace"`
		Cluster   string `json:"cluster"`
	})

	//绑定参数
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败， %v\n", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败， %v\n", err),
			"data": nil,
		})
		return
	}

	//获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用service方法，获取列表
	err = service.Servicev1.UpdateService(client, params.Content, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新Service成功",
		"data": nil,
	})

}

// CreateService 创建 Service
func (s *servicev1) CreateService(c *gin.Context) {
	//client *kubernetes.Clientset, data *ServiceCreate
	//接收参数,匿名结构体，get请求为form格式，其他请求为json格式
	var (
		serviceCreate = new(service.ServiceCreate)
	)

	//绑定参数
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := c.ShouldBindJSON(serviceCreate); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败， %v\n", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败， %v\n", err),
			"data": nil,
		})
		return
	}

	//获取client
	client, err := service.K8s.GetClient(serviceCreate.Cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	//调用service方法，获取列表
	err = service.Servicev1.CreateService(client, serviceCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "创建 Service 成功",
		"data": nil,
	})

}
