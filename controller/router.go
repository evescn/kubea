package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var Router router

type router struct{}

func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/testapi", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "testapi success!",
			"data": nil,
		})
	}).
		//pod操作
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pods/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pods", Pod.DeletePod).
		PUT("/api/k8s/pods", Pod.UpdatePod).
		GET("/api/k8s/pods/container", Pod.GetPodContainer).
		GET("/api/k8s/pods/log", Pod.GetPodLog).
		//deployment操作
		GET("/api/k8s/deployment", Deployment.GetDeployments).
		GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail).
		DELETE("/api/k8s/deployment", Deployment.DeleteDeployment).
		PUT("/api/k8s/deployment", Deployment.UpdateDeployment).
		PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment).
		PUT("/api/k8s/deployment/restart", Deployment.RestartDeployment).
		PUT("/api/k8s/deployment/create", Deployment.CreateDeployment)
}
