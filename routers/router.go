package routers

import (
	"GolangApiTest/config"
	_ "GolangApiTest/docs"
	"GolangApiTest/middleware/jwt"
	"GolangApiTest/routers/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(config.RunMode)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/user/login", api.GetAuthToken)
	r.GET("/user/info/get", api.GetUserInfo)
	r.POST("/user/logout", api.AuthLogOut)
	r.POST("/user/create", api.CreateUserNumber)
	apiToken := r.Group("/api/v1")
	apiToken.Use(jwt.JWT())
	{
		apiToken.GET("/token/:id", api.GetToken)
		apiToken.GET("/tokens", api.GetTokens)
		apiToken.POST("/token", api.AddToken)
		apiToken.PUT("/token", api.EditToken)
		apiToken.DELETE("/token/:id", api.DeleteToken)
		apiToken.GET("/svc/:id", api.GetSvcApi)
		apiToken.GET("/product/owner/total", api.GetMemberApiTotal)
		apiToken.GET("/product/total", api.GetProductApiTotal)
		apiToken.GET("/product/member", api.GetProjectMember)
	}
	return r
}
