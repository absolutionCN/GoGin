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

// @Summary 获取单个文章
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			data = models.GetArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
			logging.Error("文章不存在：", code)
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

// @Summary 获取多个文章
// @Produce  json
// @Param tag_id body int false "TagID"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [get]
func GetSomeArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS

		data["lists"] = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
		logging.Info("获取文章成功:", code)
	} else {
		for _, err := range valid.Errors {
			logging.Error(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

type AddArticlesForm struct {
	TagId     int    `gorm:"column:tag_id" json:"tag_id" valid:"Required"`
	Title     string `gorm:"column:title" json:"title" valid:"Required"`
	Desc      string `gorm:"column:desc" json:"desc" valid:"Required"`
	Content   string `gorm:"column:content" json:"content" valid:"Required"`
	CreatedBy string `gorm:"column:created_by" json:"created_by" valid:"Required"`
	State     int    `gorm:"column:state" json:"state"`
}

// @Summary 新增文章
// @Produce  json
// @Param tag_id body int true "TagID"
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param created_by body string true "CreatedBy"
// @Param state body int true "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [post]
func AddArticles(c *gin.Context) {
	var addArticlesForm AddArticlesForm
	err := c.BindJSON(&addArticlesForm)
	if err != nil {
		logging.Error("AddArticles接口获取参数失败:", err)
	}

	valid := validation.Validation{}
	valid.Min(addArticlesForm.TagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(addArticlesForm.Title, "title").Message("标题不能为空")
	valid.Required(addArticlesForm.Desc, "desc").Message("简述不能为空")
	valid.Required(addArticlesForm.Content, "content").Message("内容不能为空")
	valid.Required(addArticlesForm.CreatedBy, "created_by").Message("创建人不能为空")
	valid.Range(addArticlesForm.State, 0, 1, "state").Message("状态只允许为0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(addArticlesForm.TagId) {
			data := make(map[string]interface{})
			data["tag_id"] = addArticlesForm.TagId
			data["title"] = addArticlesForm.Title
			data["desc"] = addArticlesForm.Desc
			data["content"] = addArticlesForm.Content
			data["created_by"] = addArticlesForm.CreatedBy
			data["state"] = addArticlesForm.State

			models.AddArticle(data)
			code = e.SUCCESS
			logging.Info("新增文章成功:", code)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
			logging.Error("新增文章标签不存在:", code)
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

type EditArticlesForm struct {
	ID         int    `gorm:"column:id" json:"id" valid:"Required"`
	TagId      int    `gorm:"column:tag_id" json:"tag_id" valid:"Required"`
	Title      string `gorm:"column:title" json:"title" valid:"Required"`
	Desc       string `gorm:"column:desc" json:"desc" valid:"Required"`
	Content    string `gorm:"column:content" json:"content" valid:"Required"`
	ModifiedBy string `gorm:"column:modified_by" json:"modified_by" valid:"Required"`
	State      int    `gorm:"column:state" json:"state"`
}

// @Summary 修改文章
// @Produce  json
// @Param id path int true "ID"
// @Param tag_id body string false "TagID"
// @Param title body string false "Title"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string true "ModifiedBy"
// @Param state body int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [put]
func EditArticles(c *gin.Context) {
	var editArticlesForm EditArticlesForm
	err := c.BindJSON(&editArticlesForm)
	if err != nil {
		logging.Error("EditArticles接口获取参数失败", err)
	}
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(editArticlesForm.ID, 1, "id").Message("ID必须大于0")
	valid.MaxSize(editArticlesForm.Title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(editArticlesForm.Desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(editArticlesForm.Content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(editArticlesForm.ModifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(editArticlesForm.ModifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(editArticlesForm.ID) {
			if models.ExistTagByID(editArticlesForm.TagId) {
				data := make(map[string]interface{})
				if editArticlesForm.TagId > 0 {
					data["tag_id"] = editArticlesForm.TagId
				}
				if editArticlesForm.Title != "" {
					data["title"] = editArticlesForm.Title
				}
				if editArticlesForm.Desc != "" {
					data["desc"] = editArticlesForm.Desc
				}
				if editArticlesForm.Content != "" {
					data["content"] = editArticlesForm.Content
				}

				data["modified_by"] = editArticlesForm.ModifiedBy

				models.EditArticle(editArticlesForm.ID, data)
				code = e.SUCCESS
				logging.Info("修改文章成功:", code)
			} else {
				code = e.ERROR_NOT_EXIST_TAG
				logging.Error("修改文章TAG不存在：", code)
			}
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
			logging.Error("修改文章，文章不存在：", code)
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// @Summary 删除文章
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticles(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			models.DeleteArticle(id)
			code = e.SUCCESS
			logging.Info("删除文章成功:", code)
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
			logging.Info("删除文章。文章不存在:", code)
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
