package routers

import (
	"GoGin/ginpackage/middleware/jwt"
	"GoGin/ginpackage/pkg/setting"
	"GoGin/ginpackage/routers/api"
	"GoGin/ginpackage/routers/api/v1"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(setting.RunMode)
	r.GET("/auth", api.GetAuth)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		// 获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		// 新建标签
		apiv1.POST("/tags", v1.AddTags)
		// 更新指定标签
		apiv1.PUT("/tags", v1.EditTag)
		// 删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		// 获取文章列表
		apiv1.GET("/articles", v1.GetSomeArticles)
		// 获取指定的文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		// 新建文章
		apiv1.POST("/articles", v1.AddArticles)
		// 更新置顶文章
		apiv1.PUT("/articles/:id", v1.EditArticles)
		// 删除置顶文章
		apiv1.DELETE("articles/:id", v1.DeleteArticles)
	}
	return r
}
