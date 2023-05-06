package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubea-demo/service"
	"net/http"
)

var StatefulSet statefulSet

type statefulSet struct{}

// GetStatefulSets 获取 StatefulSet 列表
func (s *statefulSet) GetStatefulSets(c *gin.Context) {
	// client *kubernetes.Clientset, filterName, namespace string, limit, page int
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
	data, err := service.StatefulSet.GetStatefulSets(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取StatefulSet列表成功",
		"data": data,
	})

}

// GetStatefulSetDetail 获取 StatefulSet 详情
func (s *statefulSet) GetStatefulSetDetail(c *gin.Context) {
	//client *kubernetes.Clientset, dsName, namespace string
	//接收参数,匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		StsName   string `form:"sts_name"`
		Namespace string `form:"namespace"`
		Cluster   string `form:"cluster"`
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
	data, err := service.StatefulSet.GetStatefulSetDetail(client, params.StsName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取StatefulSet详情成功",
		"data": data,
	})

}

// DeleteStatefulSet 删除 StatefulSet
func (s *statefulSet) DeleteStatefulSet(c *gin.Context) {
	//client *kubernetes.Clientset, dsName, namespace string
	//接收参数,匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		StsName   string `json:"sts_name"`
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
	err = service.StatefulSet.DeleteStatefulSet(client, params.StsName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除StatefulSet成功",
		"data": nil,
	})

}

// UpdateStatefulSet 更新 StatefulSet
func (s *statefulSet) UpdateStatefulSet(c *gin.Context) {
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
	err = service.StatefulSet.UpdateStatefulSet(client, params.Content, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新StatefulSet成功",
		"data": nil,
	})

}
