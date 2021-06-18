package api

import (
	"GoGin/config/logging"
	"GoGin/config/msgCode"
	"GoGin/config/util"
	"GoGin/models"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthToken struct {
	Username string `valid:"Required; MaxSize(50)" json:"username"`
	Password string `valid:"Required; MaxSize(50)" json:"password"`
	Role     int    `json:"role"`
}

func CreateUserNumber(c *gin.Context) {
	var authToken AuthToken
	err := c.BindJSON(&authToken)
	if err != nil {
		logging.Error(err)
	}

	valid := validation.Validation{}
	//data := make(map[string]interface{})
	ok, _ := valid.Valid(&authToken)
	code := msgCode.INVALID_PARAMS
	if ok {
		if !models.ExistAuthByName(authToken.Username) {
			code = msgCode.ERROR_EXIST_AUTH
			logging.Warn(authToken.Username, " : 用户名已存在！！！！")
		} else {
			code = msgCode.SUCCESS
			models.CreateUserNumber(authToken.Username, authToken.Password, authToken.Role)
			logging.Info("创建用户成功，账号："+authToken.Username+"  密码：", authToken.Password)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": nil,
	})
}

// @Summary 登录
// @Produce  json
// @Param username path int true "username"
// @Param password path int true "password"
// @Success 200 {string} json "{"code":200,"data":{"token":"", "name:""},"msg":"ok"}"
// @Router /user/login [post]
func GetAuthToken(c *gin.Context) {
	var authToken AuthToken
	err := c.BindJSON(&authToken)
	if err != nil {
		logging.Error(err)
	}
	valid := validation.Validation{}
	data := make(map[string]interface{})
	ok, _ := valid.Valid(&authToken)
	code := msgCode.INVALID_PARAMS
	if ok {
		auth := models.CheckAuthToken(authToken.Username, authToken.Password)
		if auth != nil {
			token, err := util.GenerateToken(authToken.Username, authToken.Password)
			if err != nil {
				code = msgCode.UN_LOGGING
				logging.Warn("GetAuthToken用户未登录", authToken)
			} else {
				data["name"] = auth.Username
				data["role"] = auth.Role
				data["token"] = token
				code = msgCode.SUCCESS
				logging.Info("GetAuthToken用户登录成功", authToken)
			}
		} else {
			code = msgCode.ERROR_NOT_EXIST_AUTH
			logging.Warn("GetAuthToken用户账号密码不存在", code, authToken)
		}

	} else {
		for _, err := range valid.Errors {
			logging.Error(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": data,
	})
}

// @Summary 获取用户信息
// @Produce  json
// @Success 200 {string} json "{"code":1000,"data":{"token":"", "name:""},"msg":"ok"}"
// @Router /user/info/get [get]
func GetUserInfo(c *gin.Context) {
	data := make(map[string]interface{})
	token := c.GetHeader("Authorization")
	claims, err := util.ParseToken(token)
	code := msgCode.SUCCESS
	if err != nil {
		code = msgCode.ERROR_NOT_EXIST_AUTH
		logging.Warn("GetUserInfo用户不存在")
	} else {
		var roles []int
		auth := models.GetUserInfoByName(claims.Username)
		data["name"] = auth.Username
		roles = append(roles, auth.Role)
		data["roles"] = roles
		data["token"] = token
		//data["avatar"] = auth.Avatar
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": data,
	})
}

func AuthLogOut(c *gin.Context) {
	code := msgCode.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": nil,
	})
}
