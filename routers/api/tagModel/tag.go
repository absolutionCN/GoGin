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
	valid.Required(&addTagForm.Name, "name").Message("名称不能为空")
	valid.MaxSize(&addTagForm.Name, 100, "name").Message("名称最长为100字符")
	valid.Required(&addTagForm.CreatedBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(&addTagForm.CreatedBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(&addTagForm.State, 0, 1, "state").Message("状态只允许0或1")

	maps["name"] = addTagForm.Name
	maps["state"] = addTagForm.State
	maps["createdBy"] = addTagForm.CreatedBy

	//ok, _ := valid.Valid(&maps)
	code := msgCode.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTag(addTagForm.Name) {
			code = msgCode.ERROR_EXIST_TAG
			logging.Warn(addTagForm.Name, ": 标签已存在！！！！！！")
		} else {
			code = msgCode.SUCCESS
			models.AddTag(maps)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": maps,
	})

}

type EditTagForm struct {
	ID         int    `gorm:"column:name" json:"id" valid:"Required"`
	Name       string `gorm:"column:name" json:"tagName" valid:"Required"`
	ModifiedBy string `gorm:"column:modified_by" json:"modifiedby"`
	State      int    `gorm:"column:state" json:"state"`
}

//修改文章标签
func EditTag(c *gin.Context) {
	var editTagForm EditTagForm
	err := c.BindJSON(&editTagForm)
	if err != nil {
		logging.Error(err)
	}

	data := make(map[string]interface{})

	valid := validation.Validation{}
	valid.Required(&editTagForm.ID, "id").Message("id不能为空")
	valid.Required(&editTagForm.Name, "name").Message("名称不能为空")
	valid.MaxSize(&editTagForm.Name, 100, "name").Message("名称最长100字符")
	valid.Required(&editTagForm.ModifiedBy, "modified_by").Message("更新人不能为空")
	valid.MaxSize(&editTagForm.ModifiedBy, 100, "modified_by").Message("创建人最长为100字符")

	code := msgCode.INVALID_PARAMS

	if !valid.HasErrors() {
		if !models.ExisTagById(editTagForm.ID) {
			code = msgCode.ERROR_NOT_EXIST_TAG
			logging.Warn(editTagForm.Name, ": 标签不存在！！！！！")
		} else {
			code = msgCode.SUCCESS
			data["modified_by"] = editTagForm.ModifiedBy
			if editTagForm.Name != "" {
				data["name"] = editTagForm.Name
			}
			models.EditTag(editTagForm.ID, data)
			logging.Info("修改标签成功，传参： ", editTagForm.ID, data)
		}

	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": data,
	})

}

//删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Query("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := msgCode.INVALID_PARAMS

	if !valid.HasErrors() {
		if !models.ExisTagById(id) {
			code = msgCode.ERROR_NOT_EXIST_TAG
			logging.Warn(id, "标签不存在")
		} else {
			code = msgCode.SUCCESS
			models.DeleteTag(id)
			logging.Info("删除标签成功")
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": nil,
	})
}
