package articleModel

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

func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Query("id")).MustInt()

	valid := validation.Validation{}

	valid.Min(id, 1, "id").Message("ID必须大于0")

	var data interface{}
	code := msgCode.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistArticleByID(id) {
			code = msgCode.ERROR_NOT_EXIST_ARTICLE
			logging.Warn(id, " :标签不存在！！！")
		} else {
			code = msgCode.SUCCESS
			data = models.GetArticle(id)
			logging.Info("标签存在： ", data)
		}
	} else {
		for _, err := range valid.Errors {
			logging.Error("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": data,
	})

}

func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})

	valid := validation.Validation{}

	var state int = -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("状态必须为0或1")
	}
	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId
		valid.Min(tagId, 1, "tag_id").Message("标签id必须大于0")
	}

	var pageNum int
	if arg := c.Query("pageNum"); arg != "" {
		pageNum = com.StrTo(c.DefaultQuery("pageNum", "0")).MustInt()
	} else {
		pageNum = 0
	}

	var pageSize int
	if arg := c.Query("pageSize"); arg != "" {
		pageSize = com.StrTo(c.DefaultQuery("pageSize", "0")).MustInt()
	} else {
		pageSize = 0
	}

	code := msgCode.INVALID_PARAMS

	if !valid.HasErrors() {
		code = msgCode.SUCCESS
		data["lists"] = models.GetArticles(util.GetPage(pageNum, pageSize), config.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
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

type ArticleForm struct {
	ID         int    `gorm:"column:id" json:"id"`
	TagId      int    `gorm:"column:tag_id" json:"tag_id"`
	Title      string `gorm:"column:title" json:"title" valid:"Required"`
	Desc       string `gorm:"column:desc" json:"desc" valid:"Required"`
	Content    string `gorm:"column:content" json:"content" valid:"Required"`
	CreatedBy  string `gorm:"column:created_by" json:"createdby"`
	ModifiedBy string `gorm:"column:modified_by" json:"modifiedby"`
	State      int    `gorm:"column:state" json:"state"`
}

func AddArticle(c *gin.Context) {

	var addArticleForm ArticleForm
	data := make(map[string]interface{})

	err := c.BindJSON(&addArticleForm)
	if err != nil {
		logging.Error("新增文章报错: ", err)
	}
	valid := validation.Validation{}
	valid.Min(&addArticleForm.TagId, 1, "tag_id").Message("标签必须大于0")
	valid.Required(&addArticleForm.Title, "title").Message("标题不能为空")
	valid.Required(&addArticleForm.Desc, "desc").Message("描述不能为空")
	valid.Required(&addArticleForm.Content, "content").Message("文章内容不能为空")
	valid.Required(&addArticleForm.CreatedBy, "created_by").Message("创建人不能为空")
	valid.Range(&addArticleForm.State, 0, 1, "state").Message("状态只允许0，1")

	data["title"] = addArticleForm.Title
	data["desc"] = addArticleForm.Desc
	data["content"] = addArticleForm.Content
	data["createdBy"] = addArticleForm.CreatedBy

	code := msgCode.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExisTagById(addArticleForm.TagId) {
			data["tag_id"] = addArticleForm.TagId
			data["state"] = addArticleForm.State
			models.AddArticle(data)
			code = msgCode.SUCCESS
			logging.Info("创建文章成功：", data)
		} else {
			code = msgCode.ERROR_NOT_EXIST_TAG
			logging.Warn("标签不存在", addArticleForm.TagId)
		}
	} else {
		for _, err := range valid.Errors {
			logging.Error("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": data,
	})
}

func EditArticle(c *gin.Context) {

	var editArticleForm ArticleForm

	err := c.BindJSON(&editArticleForm)
	if err != nil {
		logging.Error("获取editarticle传参报错，报错信息为： ", err)
	}
	data := make(map[string]interface{})
	valid := validation.Validation{}
	valid.Required(&editArticleForm.ID, "id").Message("id不能为空")
	valid.MaxSize(&editArticleForm.Title, 100, "title").Message("标题最长为100个字符")
	valid.MaxSize(&editArticleForm.Desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(&editArticleForm.Content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(&editArticleForm.ModifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(&editArticleForm.ModifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	code := msgCode.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(editArticleForm.ID) {
			if models.ExisTagById(editArticleForm.TagId) {
				if editArticleForm.TagId > 0 {
					data["tag_id"] = editArticleForm.TagId
				}
				if editArticleForm.Title != "" {
					data["title"] = editArticleForm.Title
				}
				if editArticleForm.Desc != "" {
					data["desc"] = editArticleForm.Desc
				}
				if editArticleForm.Content != "" {
					data["content"] = editArticleForm.Content
				}
				data["modified_by"] = editArticleForm.ModifiedBy
				models.EditArticle(editArticleForm.ID, data)
				code = msgCode.SUCCESS
				logging.Info("修改id：", editArticleForm.ID, "文章成功，修改数据： ", data)
			} else {
				code = msgCode.ERROR_NOT_EXIST_TAG
				logging.Warn("标签不存在： ", editArticleForm.TagId)
			}
		} else {
			code = msgCode.ERROR_NOT_EXIST_ARTICLE
			logging.Warn("文章不存在：", editArticleForm.ID)
		}
	} else {
		for _, err := range valid.Errors {
			logging.Error("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": data,
	})
}

func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Query("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	code := msgCode.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistArticleByID(id) {
			code = msgCode.ERROR_NOT_EXIST_ARTICLE
			logging.Warn(id, "文章不存在")
		} else {
			code = msgCode.SUCCESS
			models.DeleteArticle(id)
			logging.Info("删除文章成功：", id)
		}
	} else {
	}
	for _, err := range valid.Errors {
		logging.Error("删除文章报错，错误信息为： ", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msgCode.GetMsg(code),
		"data": nil,
	})
}
