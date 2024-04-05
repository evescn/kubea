package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kubea/service"
	"net/http"
)

var Event event

type event struct{}

// GetList 获取 event 列表，支持过滤、排序、分页
func (*event) GetList(ctx *gin.Context) {
	params := new(struct {
		Name    string `form:"name"`
		Cluster string `form:"cluster"`
		Page    int    `form:"page"`
		Limit   int    `form:"limit"`
	})
	if err := ctx.Bind(params); err != nil {
		zap.L().Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Event.GetList(params.Name, params.Cluster, params.Page, params.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Event列表成功",
		"data": data,
	})
}
