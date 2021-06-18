package api

import (
	"GoGin/config/logging"
	"GoGin/config/msgCode"
	"GoGin/models"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 获取业务线所有人名下接口总数
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":data,"msg":"ok"}"
// @Router /api/v1/product/owner/total [get]
func GetMemberApiTotal(c *gin.Context) {
	var data []*models.Person
	code := msgCode.SUCCESS

	data, err := models.OwnerApiNumber()
	if err != nil {
		code = msgCode.ERROR
		logging.Error("获取业务线接口出错了", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": data,
	})
}

// @Summary 获取业务线所有人名下接口总数
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":data,"msg":"ok"}"
// @Router /api/v1/product/total [get]
func GetProductApiTotal(c *gin.Context) {
	var data []*models.Product
	code := msgCode.SUCCESS

	data, err := models.GetProductSummary()
	if err != nil {
		code = msgCode.ERROR
		logging.Error("获取业务线接口出错了", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": data,
	})
}

// @Summary 获取业务线人员
// @Produce  json
// @Param prodect path int true "prodect"
// @Success 200 {string} json "{"code":1000,"data":data,"msg":"ok"}"
// @Router /api/v1/product/member [get]
func GetProjectMember(c *gin.Context) {
	prodect := c.Query("prodect")
	data := make(map[string]interface{})
	valid := validation.Validation{}
	valid.Required(prodect, "prodect").Message("业务线不能为空")
	code := msgCode.INVALID_PARAMS
	if !valid.HasErrors() {
		code = msgCode.SUCCESS
		data["members"] = models.GetGroupMembers(prodect)
		data["product"] = prodect
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": data,
	})
}
