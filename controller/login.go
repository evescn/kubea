package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubea-demo/service"
	"net/http"
)

var Login login

type login struct{}

func (*login) Auth(c *gin.Context) {
	params := new(struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	})
	if err := c.ShouldBind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Login.Auth(params.UserName, params.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	// 调用页面接口
	router, err := service.VueRouter.SetRouter(data.Role)
	//logger.Info(router)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "登录成功",
		"data": router,
		"role": data.Role,
	})
}
