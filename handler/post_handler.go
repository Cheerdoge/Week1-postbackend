package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"week1-postbackend/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService interface {
	CreatePost(ctx context.Context, title, content, author string) (model.Post, error)
	UpdatePost(ctx context.Context, id primitive.ObjectID, title, content string) error
	DeletePost(ctx context.Context, id primitive.ObjectID) error
	GetPostByID(ctx context.Context, id primitive.ObjectID) (model.Post, error)
	GetAllPosts(ctx context.Context, limit int64, aftercreatedAt *time.Time, afterID *primitive.ObjectID) (posts []model.Post, nextcreatedAt *time.Time, nextID *primitive.ObjectID, err error)
	AddCommentToPost(ctx context.Context, postID primitive.ObjectID, author string, content string) (model.Comment, error)
	DeleteComment(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID) error
	UpdateCommentContent(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID, newContent string) error
	ReplyComment(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID, author string, content string) error
	DeleteReplyComment(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID, replyID primitive.ObjectID) error
	UpdateReplyCommentContent(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID, replyID primitive.ObjectID, newContent string) error
}

type PostHandler struct {
	service PostService
}

func NewPostHandler(service PostService) *PostHandler {
	return &PostHandler{service: service}
}

type PostDTO struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Body      string       `json:"body"`
	Author    string       `json:"author"`
	Comments  []CommentDTO `json:"comments"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

type CommentDTO struct {
	ID       string       `json:"id"`
	Author   string       `json:"author"`
	Content  string       `json:"content"`
	Comments []CommentDTO `json:"comments"`

	Created time.Time `json:"createdAt"`
	Updated time.Time `json:"updatedAt"`
}

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	author := c.GetString("username")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	p, err := h.service.CreatePost(ctx, req.Title, req.Content, author)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	}
	post := PostDTO{
		ID:        p.ID.Hex(),
		Title:     p.Title,
		Body:      p.Body,
		Author:    author,
		CreatedAt: p.Created,
		UpdatedAt: p.Updated,
	}
	c.JSON(200, gin.H{"post": post})
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	postIDHex := c.Param("postId")
	oid, err := primitive.ObjectIDFromHex(postIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid postId"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	err = h.service.UpdatePost(ctx, oid, req.Title, req.Content)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	p, err := h.service.GetPostByID(ctx, oid)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"post": PostDTO{
		ID:        p.ID.Hex(),
		Title:     p.Title,
		Body:      p.Body,
		Author:    p.Author,
		CreatedAt: p.Created,
		UpdatedAt: p.Updated,
	}})
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	}
	err = h.service.DeletePost(ctx, oid)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "成功"})
}

func (h *PostHandler) GetAllPost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	var afterCreatedAt *time.Time
	if s := c.Query("afterCreatedAt"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid afterCreatedAt"})
			return
		}
		afterCreatedAt = &t
	}

	var afterID *primitive.ObjectID
	if s := c.Query("afterId"); s != "" {
		oid, err := primitive.ObjectIDFromHex(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid afterId"})
			return
		}
		afterID = &oid
	}

	posts, nextCreatedAt, nextID, err := h.service.GetAllPosts(ctx, limit, afterCreatedAt, afterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var postDots []PostDTO
	for _, post := range posts {
		p := PostDTO{
			ID:        post.ID.Hex(),
			Title:     post.Title,
			Body:      post.Body,
			Author:    post.Author,
			CreatedAt: post.Created,
			UpdatedAt: post.Updated,
		}
		postDots = append(postDots, p)
	}

	c.JSON(200, gin.H{
		"posts":              postDots,
		"nextAfterCreatedAt": nextCreatedAt.Format(time.RFC3339),
		"nextAfterId":        nextID.Hex(),
	})
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	post, err := h.service.GetPostByID(ctx, oid)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	var commentDTOs []CommentDTO
	for _, comment := range post.Comments {
		var replyDTOs []CommentDTO
		for _, reply := range comment.Comments {
			r := CommentDTO{
				ID:      reply.ID.Hex(),
				Author:  reply.Author,
				Content: reply.Content,
				Created: reply.Created,
				Updated: reply.Updated,
			}
			replyDTOs = append(replyDTOs, r)
		}
		p := CommentDTO{
			ID:       comment.ID.Hex(),
			Author:   comment.Author,
			Content:  comment.Content,
			Comments: replyDTOs,
			Created:  comment.Created,
			Updated:  comment.Updated,
		}
		commentDTOs = append(commentDTOs, p)
	}
	c.JSON(200, gin.H{"post": PostDTO{
		ID:        post.ID.Hex(),
		Title:     post.Title,
		Body:      post.Body,
		Author:    post.Author,
		Comments:  commentDTOs,
		CreatedAt: post.Created,
		UpdatedAt: post.Updated,
	}})
}

type AddCommentToPostRequest struct {
	Content string `json:"content"`
}

func (h *PostHandler) AddCommentToPost(c *gin.Context) {
	author := c.GetString("username")
	var req AddCommentToPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	comment, err := h.service.AddCommentToPost(ctx, oid, author, req.Content)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"comment": CommentDTO{
		ID:      comment.ID.Hex(),
		Author:  comment.Author,
		Content: comment.Content,
		Created: comment.Created,
		Updated: comment.Updated,
	}})
}

func (h *PostHandler) DeleteComment(c *gin.Context) {
	poid, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	coid, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.service.DeleteComment(ctx, poid, coid)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "成功"})
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

func (h *PostHandler) UpdateCommentContent(c *gin.Context) {
	poid, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	coid, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.service.UpdateCommentContent(ctx, poid, coid, req.Content)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "成功"})
}

type RemoveCommentRequest struct {
	Content string `json:"content"`
}

func (h *PostHandler) ReplyComment(c *gin.Context) {
	author := c.GetString("username")
	poid, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coid, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req RemoveCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.service.ReplyComment(ctx, poid, coid, author, req.Content)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "回复成功"})
}

func (h *PostHandler) DeleteReplyComment(c *gin.Context) {
	poid, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coid, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roid, err := primitive.ObjectIDFromHex(c.Param("replyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.service.DeleteReplyComment(ctx, poid, coid, roid)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "删除成功"})
}

type UpdateReplyCommentRequest struct {
	Content string `json:"content"`
}

func (h *PostHandler) UpdateReplyCommentContent(c *gin.Context) {
	poid, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coid, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roid, err := primitive.ObjectIDFromHex(c.Param("replyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req UpdateReplyCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.service.UpdateReplyCommentContent(ctx, poid, coid, roid, req.Content)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "修改成功"})

}
