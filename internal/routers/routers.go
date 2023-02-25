package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/taoruicheng/blog-service/global"
	"github.com/taoruicheng/blog-service/internal/middleware"
	v1 "github.com/taoruicheng/blog-service/internal/routers/api/v1"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Tracing())
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}

	r.Use(middleware.RateLimiter(global.MethodLimiterSetting))
	r.Use(middleware.Translations())

	tag := v1.NewTag()

	apiv1 := r.Group("/api/v1")

	// 创建标签
	apiv1.POST("/tags", tag.Create)
	// 删除指定标签
	apiv1.DELETE("/tags/:id", tag.Delete)
	// 更新指定标签
	apiv1.PUT("/tags/:id", tag.Update)
	// 获取标签列表
	apiv1.GET("/tags", tag.List)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
