package handler

import (
	"context"
	"time"
	"week1-postbackend/model"

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
