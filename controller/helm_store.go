package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kubea/model"
	"kubea/service"
	"net/http"
)

var HelmStore helmStore

type helmStore struct{}

// ListReleases 已安装的release列表
func (*helmStore) ListReleases(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Cluster    string `form:"cluster"`
	})
	if err := ctx.Bind(params); err != nil {
		zap.L().Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	actionConfig, err := service.HelmConfig.GetAction(params.Cluster, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.HelmStore.ListReleases(actionConfig, params.FilterName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Release列表成功",
		"data": data,
	})
}

// DetailRelease release详情
func (*helmStore) DetailRelease(ctx *gin.Context) {
	params := new(struct {
		Release   string `form:"release"`
		Namespace string `form:"namespace"`
		Cluster   string `form:"cluster"`
	})
	if err := ctx.Bind(params); err != nil {
		zap.L().Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	actionConfig, err := service.HelmConfig.GetAction(params.Cluster, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.HelmStore.DetailRelease(actionConfig, params.Release)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Release详情成功",
		"data": data,
	})
}

// InstallRelease release安装
func (*helmStore) InstallRelease(ctx *gin.Context) {
	params := new(struct {
		Release   string `json:"release"`
		Chart     string `json:"chart"`
		Namespace string `json:"namespace"`
		Cluster   string `json:"cluster"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	actionConfig, err := service.HelmConfig.GetAction(params.Cluster, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.HelmStore.InstallRelease(actionConfig, params.Release, params.Chart, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "安装Release成功",
		"data": nil,
	})
}

// UninstallRelease release卸载
func (*helmStore) UninstallRelease(ctx *gin.Context) {
	params := new(struct {
		Release   string `json:"release"`
		Namespace string `json:"namespace"`
		Cluster   string `json:"cluster"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	actionConfig, err := service.HelmConfig.GetAction(params.Cluster, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.HelmStore.UninstallRelease(actionConfig, params.Release)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "卸载Release成功",
		"data": nil,
	})
}

// UploadChartFile chart文件上传
func (*helmStore) UploadChartFile(ctx *gin.Context) {
	//ctx.Request.FormFile可以直接获取到file multipart.File, header *multipart.FileHeader
	file, header, err := ctx.Request.FormFile("chart")
	if err != nil {
		zap.L().Error("获取上传信息失败, " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.HelmStore.UploadChartFile(file, header)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "上传Chart文件成功",
		"data": nil,
	})
}

// DeleteChartFile chart文件删除
func (*helmStore) DeleteChartFile(ctx *gin.Context) {
	params := new(struct {
		Chart string `json:"chart"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.HelmStore.DeleteChartFile(params.Chart)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除Chart文件成功",
		"data": nil,
	})
}

// ListCharts chart列表
func (*helmStore) ListCharts(ctx *gin.Context) {
	params := new(struct {
		Name  string `form:"name"`
		Page  int    `form:"page"`
		Limit int    `form:"limit"`
	})
	if err := ctx.Bind(params); err != nil {
		zap.L().Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.HelmStore.ListCharts(params.Name, params.Page, params.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Chart列表成功",
		"data": data,
	})
}

// AddChart chart新增
func (*helmStore) AddChart(ctx *gin.Context) {
	params := new(model.Chart)
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.HelmStore.AddChart(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "新增Chart成功",
		"data": nil,
	})
}

// UpdateChart chart更新
func (*helmStore) UpdateChart(ctx *gin.Context) {
	params := new(model.Chart)
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.HelmStore.UpdateChart(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新Chart成功",
		"data": nil,
	})
}

// DeleteChart chart删除
func (*helmStore) DeleteChart(ctx *gin.Context) {
	params := new(model.Chart)
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.HelmStore.DeleteChart(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除Chart成功",
		"data": nil,
	})
}
