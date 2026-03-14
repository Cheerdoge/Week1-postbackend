package router

import (
	"week1-postbackend/handler"
	"week1-postbackend/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth *handler.AuthHandler
	// 你后续可以加 PostHandler
}

func NewRouter(h Handlers) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
	}

	// 需要登录的接口示例
	protected := api.Group("")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/me", func(c *gin.Context) {
			// demo：验证 token 可用
			c.JSON(200, gin.H{
				"userID":   c.MustGet("userID"),
				"email":    c.MustGet("email"),
				"username": c.MustGet("username"),
			})
		})
	}

	return r
}
