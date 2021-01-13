package api

import (
	"GoGin/ginpackage/models"
	"GoGin/ginpackage/pkg/logging"
	"GoGin/ginpackage/pkg/setting/e"
	"GoGin/ginpackage/pkg/util"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary 获取用户
// @Produce  json
// @Param username query string true "userName"
// @Param password query string true "password"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
				logging.Error("getAuth校验token不通过:", code)
			} else {
				data["token"] = token

				code = e.SUCCESS
				logging.Info("getAuth调用成功:", code)
			}
		} else {
			code = e.ERROR_AUTH
			logging.Info("用户权限错误：", code)
		}

	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
