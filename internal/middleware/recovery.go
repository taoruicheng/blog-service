package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/taoruicheng/blog-service/global"
	"github.com/taoruicheng/blog-service/pkg/app"
	"github.com/taoruicheng/blog-service/pkg/errcode"
)

func Recovery() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logger.WithCallersFrames().Errorf(c, "panic recover err: %v", err)

				app.NewResponse(c).ToErrorResponse(errcode.ServerError)
				c.Abort()
			}
		}()
		c.Next()
	}
}
