package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubea-demo/service"
	"net/http"
)

var Pvc pvc

type pvc struct{}

// GetPvcs 获取 PVC 列表
func (p *pvc) GetPvcs(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
		Cluster    string `form:"cluster"`
	})

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

	//调用service方法，
	data, err := service.Pvc.GetPvcs(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取PVC列表成功",
		"data": data,
	})

}

// GetPvcDetail 获取 PVC 详情
func (p *pvc) GetPvcDetail(c *gin.Context) {
	params := new(struct {
		PvcName   string `form:"pvc_name"`
		Namespace string `form:"namespace"`
		Cluster   string `form:"cluster"`
	})

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
	data, err := service.Pvc.GetPvcDetail(client, params.PvcName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取PVC详情成功",
		"data": data,
	})

}

// DeletePvc 删除 PVC
func (p *pvc) DeletePvc(c *gin.Context) {
	params := new(struct {
		PvcName   string `json:"pvc_name"`
		Namespace string `json:"namespace"`
		Cluster   string `json:"cluster"`
	})

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
	err = service.Pvc.DeletePvc(client, params.PvcName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除PVC成功",
		"data": nil,
	})

}

// UpdatePvc 更新 Pvc
func (p *pvc) UpdatePvc(c *gin.Context) {
	params := new(struct {
		Content   string `json:"content"`
		Namespace string `json:"namespace"`
		Cluster   string `json:"cluster"`
	})

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

	//调用service方法
	err = service.Pvc.UpdatePvc(client, params.Content, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新PVC成功",
		"data": nil,
	})

}
