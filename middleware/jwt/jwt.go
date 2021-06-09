package jwt

import (
	"GolangApiTest/config/logging"
	"GolangApiTest/config/msgCode"
	"GolangApiTest/config/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = msgCode.SUCCESS
		token := c.GetHeader("Authorization")
		if token == "" {
			code = msgCode.ERROR_NOT_EXIST_AUTH
			logging.Error("用户token不存在")
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = msgCode.ERROR_AUTH_TOKEN_CHECK_FAIL
				logging.Error("用户token校验失败")
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = msgCode.ERROR_AUTH_TOKEN_CHECK_TIMEOUT
				logging.Error("用户token已过期")
			}
		}
		if code != msgCode.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  msgCode.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}
		c.Next()
	}
}
