package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kubea/service"
	"net/http"
)

var Secret secret

type secret struct{}

// GetSecrets 获取 Secret 列表
func (p *secret) GetSecrets(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
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
	data, err := service.Secret.GetSecrets(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Secret列表成功",
		"data": data,
	})

}

// GetSecretDetail 获取 Secret 详情
func (p *secret) GetSecretDetail(c *gin.Context) {
	params := new(struct {
		SecretName string `form:"secret_name"`
		Namespace  string `form:"namespace"`
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
	data, err := service.Secret.GetSecretDetail(client, params.SecretName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Secret详情成功",
		"data": data,
	})

}

// DeleteSecret 删除 Secret
func (p *secret) DeleteSecret(c *gin.Context) {
	params := new(struct {
		SecretName string `json:"secret_name"`
		Namespace  string `json:"namespace"`
		Cluster    string `json:"cluster"`
	})

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
	err = service.Secret.DeleteSecret(client, params.SecretName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除Secret成功",
		"data": nil,
	})

}

// UpdateSecret 更新 Secret
func (p *secret) UpdateSecret(c *gin.Context) {
	params := new(struct {
		Content   string `json:"content"`
		Namespace string `json:"namespace"`
		Cluster   string `json:"cluster"`
	})

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
	err = service.Secret.UpdateSecret(client, params.Content, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新Secret成功",
		"data": nil,
	})

}
