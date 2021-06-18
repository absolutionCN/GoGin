package api

import (
	"GoGin/config/logging"
	"GoGin/config/msgCode"
	"GoGin/config/util"
	"GoGin/models"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

// @Summary 获取单个token
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/token/{id} [get]
func GetToken(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := msgCode.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors() {
		if models.EditTokenId(id) {
			data = models.GetToken(id)
			code = msgCode.SUCCESS
			logging.Info("GetToken接口获取到的id为：", id)
		} else {
			code = msgCode.ERROR_NOT_EXIST_TOKENID
			logging.Warn("GetTokenw未获取到id", msgCode.ERROR_NOT_EXIST_TOKENID)
		}
	} else {
		for _, err := range valid.Errors {
			//log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
			logging.Error(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": data,
	})
}

// @Summary 获取token列表
// @Produce  json
// @Param state path int false "state"
// @Param pageNum path int false "pageNum"
// @Param pageSize path int false "pageSize"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tokens [get]
func GetTokens(c *gin.Context) {
	maps := make(map[string]interface{})
	data := make(map[string]interface{}) //page := c.Query("page")

	//pageNum := com.StrTo(c.Query("pageNum")).MustInt()
	//pageSize := com.StrTo(c.Query("pageSize")).MustInt()
	valid := validation.Validation{}

	//var state int = -1
	var state int
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(c.DefaultQuery("state", "0")).MustInt()
		//maps["state"] = state
	}

	var pageNum int
	if arg := c.Query("pageNum"); arg != "" {
		pageNum = com.StrTo(c.DefaultQuery("pageNum", "0")).MustInt()
	} else {
		pageNum = 0
	}

	var pageSize int
	if arg := c.Query("pageSize"); arg != "" {
		pageSize = com.StrTo(c.DefaultQuery("pageSize", "10")).MustInt()
	} else {
		pageSize = 10
	}

	code := msgCode.INVALID_PARAMS
	if !valid.HasErrors() {
		code = msgCode.SUCCESS
		data["lists"] = models.GetTokens(util.GetPage(pageNum, pageSize), pageSize, state)
		data["total"] = models.GetTokenTotal(maps)
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

type AddTokenForm struct {
	ProjectName string `gorm:"column:prodect" json:"prodect" valid:"Required"`
	SvcName     string `gorm:"servername" json:"servername" valid:"Required"`
	Sid         int    `gorm:"sid" json:"sid" valid:"Required"`
	Token       string `gorm:"token" json:"token" valid:"Required"`
	Owner       string `gorm:"owner" json:"owner" valid:"Required"`
}

// @Summary 新增token
// @Produce  json
// @Param prodect body string true "prodect"
// @Param servername body string true "servername"
// @Param token body string true "token"
// @Param sid body int true "sid"
// @Param owner body string true "owner"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/token [post]
func AddToken(c *gin.Context) {
	var formToken AddTokenForm
	err := c.BindJSON(&formToken)
	if err != nil {
		logging.Error(err)
	}
	valid := validation.Validation{}
	valid.Required(&formToken.ProjectName, "prodect").Message("prodect不能为空")
	valid.MaxSize(&formToken.ProjectName, 100, "prodect").Message("prodect最长为100字符")
	valid.Required(&formToken.SvcName, "svcname").Message("svcname不能为空")
	valid.MaxSize(&formToken.SvcName, 100, "svcname").Message("svcname最长为100字符")
	valid.Required(&formToken.Sid, "sid").Message("sid不能为空")
	valid.Required(&formToken.Token, "token").Message("token不能为空")
	valid.MaxSize(&formToken.Token, 100, "token").Message("token最长为255字符")
	valid.Required(&formToken.Owner, "owner").Message("owner不能为空")

	code := msgCode.INVALID_PARAMS

	if !valid.HasErrors() {
		if !models.ExistToken(formToken.Token) {
			tokenData := make(map[string]interface{})
			tokenData["prodect"] = formToken.ProjectName
			tokenData["servername"] = formToken.SvcName
			tokenData["sid"] = formToken.Sid
			tokenData["token"] = formToken.Token
			models.AddToken(tokenData)
			code = msgCode.SUCCESS
			logging.Info("AddToken接口添加数据：", tokenData)
		} else {
			code = msgCode.ERROR_EXIST_TOKEN
			logging.Warn("AddToken传入token为：", formToken.Token)
		}
	} else {
		for _, err := range valid.Errors {
			//log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

type EditTokenForm struct {
	ID          int    `gorm:"primary_key" json:"id"`
	ProjectName string `gorm:"column:prodect" json:"prodect" valid:"Required"`
	SvcName     string `gorm:"servername" json:"servername" valid:"Required"`
	Sid         int    `gorm:"sid" json:"sid" valid:"Required"`
	Token       string `gorm:"token" json:"token" valid:"Required"`
	Coverage    int    `gorm:"coverage" json:"coverage" valid:"Required"`
	Owner       string `gorm:"owner" json:"owner" valid:"Required"`
}

// @Summary 修改token
// @Produce  json
// @Param id path int true "ID"
// @Param prodect body string true "prodect"
// @Param servername body string true "servername"
// @Param token body string true "token"
// @Param sid body int true "sid"
// @Param coverage body int true "coverage"
// @Param owner body string true "owner"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/token [put]
func EditToken(c *gin.Context) {
	var editToken EditTokenForm
	err := c.BindJSON(&editToken)
	if err != nil {
		logging.Error(err)
	}

	valid := validation.Validation{}
	valid.Min(&editToken.ID, 1, "id").Message("ID必须大于0")
	valid.Required(&editToken.ProjectName, "prodect").Message("prodect不能为空")
	valid.MaxSize(&editToken.ProjectName, 100, "prodect").Message("prodect最长为100字符")
	valid.Required(&editToken.SvcName, "svcname").Message("svcname不能为空")
	valid.MaxSize(&editToken.SvcName, 100, "svcname").Message("svcname最长为100字符")
	valid.Required(&editToken.Sid, "sid").Message("sid不能为空")
	valid.Required(&editToken.Token, "token").Message("token不能为空")
	valid.MaxSize(&editToken.Token, 100, "token").Message("token最长为255字符")
	valid.MaxSize(&editToken.Owner, 100, "owner").Message("owner最长为100字符")
	valid.Required(&editToken.Owner, "owner").Message("owner不能为空")

	code := msgCode.INVALID_PARAMS

	if !valid.HasErrors() {
		if models.EditTokenId(editToken.ID) {
			data := make(map[string]interface{})
			//if editToken.ID  > 0 {
			//	data["id"] = editToken.ID
			//}
			if editToken.Token != "" {
				data["token"] = editToken.Token
			}
			if editToken.ProjectName != "" {
				data["prodect"] = editToken.ProjectName
			}
			if editToken.SvcName != "" {
				data["servername"] = editToken.SvcName
			}
			if editToken.Sid > 0 {
				data["sid"] = editToken.Sid
			}
			if editToken.Coverage > 0 {
				data["coverage"] = editToken.Coverage
			}
			if editToken.Owner != "" {
				data["owner"] = editToken.Owner
			}
			models.EditToken(editToken.ID, data)
			code = msgCode.SUCCESS
			logging.Info("EditToken接口传参", data)
		} else {
			code = msgCode.ERROR_NOT_EXIST_TOKENID
			logging.Warn("EditToken接token不存在", code)
		}
	} else {
		for _, err := range valid.Errors {
			//log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": make(map[string]string),
	})
}

// @Summary 删除token
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/token [delete]
func DeleteToken(c *gin.Context) {

	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := msgCode.INVALID_PARAMS

	if !valid.HasErrors() {
		if models.EditTokenId(id) {
			models.DeleteToken(id)
			code = msgCode.SUCCESS
			logging.Info("DeleteToken传入参数", id)
		} else {
			code = msgCode.ERROR_NOT_EXIST_TOKENID
			logging.Warn("DeleteToken传入token不存在", code)
		}
	} else {
		for _, err := range valid.Errors {
			//log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{

		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": make(map[string]string),
	})
}

// @Summary 获取服务所有接口
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/svc/{id} [get]
func GetSvcApi(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := msgCode.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors() {
		rs := models.SyncSvcApi(id)
		if rs == msgCode.TaskError {
			code = msgCode.ERROR_SYNC_YAPI
			logging.Error(msgCode.ERROR_SYNC_YAPI)
		}
		if rs == msgCode.TaskRunning {
			code = msgCode.ERROR_SYNC_RUNNING
			logging.Error(msgCode.ERROR_SYNC_RUNNING)
		}
		if rs == msgCode.TaskSuccess {
			data = struct{}{}
			code = msgCode.SUCCESS
			logging.Info("开始同步yapi接口信息：", id)
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
