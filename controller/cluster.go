package controller

import (
	"github.com/gin-gonic/gin"
	"kubea/service"
	"net/http"
	"sort"
)

var Cluster cluster

type cluster struct{}

func (*cluster) GetClusters(c *gin.Context) {
	list := make([]string, 0)
	for key := range service.K8s.ClientMap {
		list = append(list, key)
	}
	sort.Strings(list)
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取集群信息成功",
		"data": list,
	})
}
