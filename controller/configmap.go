package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kubea/service"
	"net/http"
)

var ConfigMap configmap

type configmap struct{}

// GetConfigMaps 获取 ConfigMap 列表
func (p *configmap) GetConfigMaps(c *gin.Context) {
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

	//调用service方法，获取列表
	data, err := service.ConfigMap.GetConfigMaps(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取ConfigMap列表成功",
		"data": data,
	})

}

// GetConfigMapDetail 获取 ConfigMap 详情
func (p *configmap) GetConfigMapDetail(c *gin.Context) {
	params := new(struct {
		CmName    string `form:"cm_name"`
		Namespace string `form:"namespace"`
		Cluster   string `form:"cluster"`
	})

	//绑定参数
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := c.Bind(params); err != nil {
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

	//调用service方法，获取列表
	data, err := service.ConfigMap.GetConfigMapDetail(client, params.CmName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取ConfigMap详情成功",
		"data": data,
	})

}

// DeleteConfigMap 删除 ConfigMap
func (p *configmap) DeleteConfigMap(c *gin.Context) {
	//client *kubernetes.Clientset, podName, namespace string
	//接收参数,匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		CmName    string `json:"cm_name"`
		Namespace string `json:"namespace"`
		Cluster   string `json:"cluster"`
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

	//调用service方法，获取列表
	err = service.ConfigMap.DeleteConfigMap(client, params.CmName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除ConfigMap成功",
		"data": nil,
	})

}

// UpdateConfigMap 更新 ConfigMap
func (p *configmap) UpdateConfigMap(c *gin.Context) {
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

	//调用service方法
	err = service.ConfigMap.UpdateConfigMap(client, params.Content, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新ConfigMap成功",
		"data": nil,
	})

}
