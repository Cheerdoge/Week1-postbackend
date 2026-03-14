package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 支持 file:// (Origin 会是 "null")；同时也允许常见本地开发端口
		allowed := map[string]bool{
			"null":                  true,
			"http://localhost:5500": true,
			"http://127.0.0.1:5500": true,
			"http://localhost:5173": true,
			"http://127.0.0.1:5173": true,
		}

		if allowed[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Max-Age", "86400")
		}

		// 预检请求必须直接返回，否则浏览器会报 Failed to fetch
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
