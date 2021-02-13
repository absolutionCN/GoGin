package v1

import (
	"GoGin/ginpackage/models"
	"GoGin/ginpackage/pkg/logging"
	"GoGin/ginpackage/pkg/setting"
	"GoGin/ginpackage/pkg/setting/e"
	"GoGin/ginpackage/pkg/util"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// @Summary 获取标签
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}
	code := e.SUCCESS

	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

type AddTagsForm struct {
	Name      string `gorm:"column:name" json:"name" valid:"Required"`
	CreatedBy string `gorm:"column:created_by" json:"created_by" valid:"Required"`
	State     int    `gorm:"column:state" json:"state"`
}

// @Summary 新增文章标签
// @Produce  json
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTags(c *gin.Context) {
	var addTagsForm AddTagsForm
	err := c.BindJSON(&addTagsForm)
	if err != nil {
		logging.Error("AddTags获取参数失败：", err)
	}
	valid := validation.Validation{}
	valid.Required(addTagsForm.Name, "name").Message("名称不能为空")
	valid.MaxSize(addTagsForm.Name, 100, "name").Message("名称最长为100字符")
	valid.Required(addTagsForm.CreatedBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(addTagsForm.CreatedBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(addTagsForm.State, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByName(addTagsForm.Name) {
			code = e.SUCCESS
			models.AddTag(addTagsForm.Name, addTagsForm.State, addTagsForm.CreatedBy)
			logging.Info("新增tag成功")
		} else {
			code = e.ERROR_EXIST_TAG
			logging.Error("tag已存在：", code)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

type EditTagForm struct {
	ID         int    `gorm:"column:id" json:"id" valid:"Required"`
	Name       string `gorm:"column:name" json:"name" valid:"Required"`
	ModifiedBy string `gorm:"column:modified_by" json:"modified_by" valid:"Required"`
}

// @Summary 修改文章标签
// @Produce  json
// @Param id path int true "ID"
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param modified_by body string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	// 修改文章标签
	var editTagForm EditTagForm
	err := c.BindJSON(&editTagForm)
	if err != nil {
		logging.Error("EditTag接口获取参数失败:", err)
	}
	vaild := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		vaild.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	vaild.Required(editTagForm.ID, "id").Message("ID不能为空")
	vaild.Required(editTagForm.ModifiedBy, "modified_by").Message("修改人不能为空")
	vaild.MaxSize(editTagForm.ModifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	vaild.MaxSize(editTagForm.Name, 100, "name").Message("名字最长为100字符")

	code := e.INVALID_PARAMS
	if !vaild.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(editTagForm.ID) {
			data := make(map[string]interface{})
			data["modified_by"] = editTagForm.ModifiedBy
			if editTagForm.Name != "" {
				data["name"] = editTagForm.Name
			}
			if state != -1 {
				data["state"] = state
			}

			models.EditTag(editTagForm.ID, data)
			logging.Info("修改后的标签数据为：", data)
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
			logging.Error("文章不存在：", code)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// @Summary 删除文章标签
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	//删除文章标签
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			models.DeleteTag(id)
			logging.Info("删除文章标签成功, 标签id:", id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
			logging.Error("tag不存在：", code)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
