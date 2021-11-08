package server

import (
	"api-web/api"
	"api-web/middleware"
	"api-web/model"
	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	// 配置权限
	var authMiddleware = middleware.GinJWTMiddlewareInit()

	// 初始化
	gin.ForceConsoleColor()
	r := gin.Default()

	// 配置跨域
	r.Use(middleware.Cors())

	// 路由
	v := r.Group("/api")
	{
		v.GET("ping", api.Ping)

		// 登录保护的路由
		auth := v.Group("/v1")
		auth.Use(authMiddleware.MiddlewareFunc())
		{
			// 挂载文件
			auth.POST("upload", model.UploadFile)
			auth.GET("download/:hash", model.DownloadFile)
			// 挂载Graphql
			auth.GET("/graphql", middleware.GraphqlHandler())
			auth.POST("/graphql", middleware.GraphqlHandler())
		}
	}
	return r
}
