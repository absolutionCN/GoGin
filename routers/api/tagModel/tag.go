package tagModel

import (
	"GoGin/config"
	"GoGin/config/logging"
	"GoGin/config/msgCode"
	"GoGin/config/util"
	"GoGin/models"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

//获取多个文章标签
func GetTags(c *gin.Context) {

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	valid := validation.Validation{}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
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
		data["lists"] = models.GetTags(util.GetPage(pageNum, pageSize), config.PageSize, maps)
		data["total"] = models.GetTagTotal(maps)
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

type AddTagForm struct {
	Name      string `gorm:"column:name" json:"tagName" valid:"Required"`
	State     int    `gorm:"column:state" json:"state" valid:"Required"`
	CreatedBy string `gorm:"column:created_by" json:"createdBy" valid:"Required"`
}

//新增文章标签
func AddTag(c *gin.Context) {
	var addTagForm AddTagForm

	maps := make(map[string]interface{})

	err := c.BindJSON(&addTagForm)
	if err != nil {
		logging.Error(err)
	}
	valid := validation.Validation{}
	valid.Required(&addTagForm.Name, "name").Message("标签名称不能为空")
	valid.MaxSize(&addTagForm.Name, 100, "name").Message("标签名称最长为100个字符")
	valid.Required(&addTagForm.State, "state").Message("标签状态不能为空")
	valid.MaxSize(&addTagForm.State, 3, "state").Message("标签状态最长为3个字符")
	valid.Required(&addTagForm.CreatedBy, "createdBy").Message("创建人不能为空")
	valid.MaxSize(&addTagForm.CreatedBy, 100, "createdBy").Message("标签创建人最长为100个字符")

	maps["name"] = addTagForm.Name
	maps["state"] = addTagForm.State
	maps["createdBy"] = addTagForm.CreatedBy

	//ok, _ := valid.Valid(&maps)
	code := msgCode.INVALID_PARAMS
	if models.ExistTag(addTagForm.Name) {
		code = msgCode.ERROR_EXIST_TAG
		logging.Warn(addTagForm.Name, ": 标签已存在！！！！！！")
	} else {
		code = msgCode.SUCCESS
		models.AddTag(maps)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": nil,
	})

}

//修改文章标签
func EditTag(c *gin.Context) {

}

//删除文章标签
func DeleteTag(c *gin.Context) {

}
