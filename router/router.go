package router

import (
	"week1-postbackend/handler"
	"week1-postbackend/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth *handler.AuthHandler
	Post *handler.PostHandler
}

func NewRouter(h Handlers) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	api := r.Group("/api/v1")

	// auth：不需要登录
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
	}

	// protected：全部需要登录
	protected := api.Group("")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"userID":   c.MustGet("userID"),
				"email":    c.MustGet("email"),
				"username": c.MustGet("username"),
			})
		})

		posts := protected.Group("/posts")
		{
			// post CRUD
			posts.GET("", h.Post.GetAllPost)          // GET /api/v1/posts?limit=&afterCreatedAt=&afterId=
			posts.GET("/:postId", h.Post.GetPostByID) // GET /api/v1/posts/{postId}
			posts.POST("", h.Post.CreatePost)         // POST /api/v1/posts
			posts.PATCH("/:postId", h.Post.UpdatePost)
			posts.DELETE("/:postId", h.Post.DeletePost)

			// 一级评论
			posts.POST("/:postId/comments", h.Post.AddCommentToPost)
			posts.PATCH("/:postId/comments/:commentId", h.Post.UpdateCommentContent)
			posts.DELETE("/:postId/comments/:commentId", h.Post.DeleteComment)

			// 二级评论（回复）
			posts.POST("/:postId/comments/:commentId/replies", h.Post.ReplyComment)
			posts.PATCH("/:postId/comments/:commentId/replies/:replyId", h.Post.UpdateReplyCommentContent)
			posts.DELETE("/:postId/comments/:commentId/replies/:replyId", h.Post.DeleteReplyComment)
		}
	}

	return r
}
