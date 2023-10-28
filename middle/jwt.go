package middle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kubea-demo/utils"
	"net/http"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//if len(c.Request.URL.String()) >= 10 {
		fmt.Println(c.Request.URL.String())
		if len(c.Request.URL.String()) >= 10 && c.Request.URL.String()[0:10] == "/api/login" {
			c.Next()
		} else if len(c.Request.URL.String()) >= 12 && c.Request.URL.String()[0:12] == "/api/app/get" {
			c.Next()
		} else if len(c.Request.URL.String()) >= 15 && c.Request.URL.String()[0:15] == "/api/deploy/add" {
			c.Next()
		} else if len(c.Request.URL.String()) >= 23 && c.Request.URL.String()[0:23] == "/api/deploy/jenkniscicd" {
			c.Next()
		} else if len(c.Request.URL.String()) >= 22 && c.Request.URL.String()[0:22] == "/api/deploy/updatecicd" {
			c.Next()
		} else {
			token := c.Request.Header.Get("Authorization")
			if token == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg":  "请求未携带token，无权限访问",
					"data": nil,
				})
				c.Abort()
				return
			}

			claims, err := utils.JWTToken.ParseToken(token)
			if err != nil {
				//token延期错误
				if err.Error() == "TokenExpired" {
					c.JSON(http.StatusBadRequest, gin.H{
						"msg":  "授权已过期",
						"data": nil,
					})
					c.Abort()
					return
				}
				//其他解析错误
				c.JSON(http.StatusBadRequest, gin.H{
					"msg":  err.Error(),
					"data": nil,
				})
				c.Abort()
				return
			}

			// 继续交由下一个路由处理,并将解析出的信息传递下去
			c.Set("claims", claims)
			c.Next()
		}
	}
}
