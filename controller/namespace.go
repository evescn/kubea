package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kubea/service"
	"net/http"
)

var Namespace namespace

type namespace struct{}

// GetNamespaces 获取 Namespace 列表
func (n *namespace) GetNamespaces(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
		Cluster    string `form:"cluster"`
	})

	if err := c.Bind(params); err != nil {
		zap.L().Error(fmt.Sprintf("绑定参数失败， %v\n", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败， %v\n", err),
			"data": nil,
		})
		return
	}

	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Namespace.GetNamespaces(client, params.FilterName, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Namespace列表成功",
		"data": data,
	})
}

// GetNamespaceDetail 获取 Namespace 列表
func (n *namespace) GetNamespaceDetail(c *gin.Context) {
	params := new(struct {
		NamespaceName string `form:"namespace_name"`
		Cluster       string `form:"cluster"`
	})

	if err := c.Bind(params); err != nil {
		zap.L().Error(fmt.Sprintf("绑定参数失败， %v\n", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败， %v\n", err),
			"data": nil,
		})
		return
	}

	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Namespace.GetNamespaceDetail(client, params.NamespaceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Namespace详情成功",
		"data": data,
	})
}

// DeleteNamespace 删除 Namespace
func (n *namespace) DeleteNamespace(c *gin.Context) {
	//client *kubernetes.Clientset, podName, namespace string
	//接收参数,匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		NamespaceName string `json:"namespace_name"`
		Cluster       string `json:"cluster"`
	})

	//绑定参数
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := c.ShouldBindJSON(params); err != nil {
		zap.L().Error(fmt.Sprintf("绑定参数失败， %v\n", err))
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

	//调用service方法，删除
	err = service.Namespace.DeleteNamespace(client, params.NamespaceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除Namespace成功",
		"data": nil,
	})

}
