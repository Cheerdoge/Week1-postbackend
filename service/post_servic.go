package service

import (
	"context"
	"errors"
	"time"
	"week1-postbackend/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostRepository interface {
	CreatePost(ctx context.Context, title, content, author string) (model.Post, error)
	UpdatePost(ctx context.Context, id primitive.ObjectID, title, content string) error
	DeletePost(ctx context.Context, id primitive.ObjectID) error
	GetPostByID(ctx context.Context, id primitive.ObjectID) (model.Post, error)
	ListByCursor(ctx context.Context, limit int64, afterCreated *time.Time, afterID *primitive.ObjectID) (posts []model.Post, nextCreatedAt *time.Time, nextID *primitive.ObjectID, err error)
	AddCommentToPost(ctx context.Context, postID primitive.ObjectID, author string, content string) (model.Comment, error)
	DeleteComment(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID) error
	UpdateComment(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID, newContent string) error
	ReplyComment(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID, author string, content string) error
	DeleteReplyComment(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID, replyID primitive.ObjectID) error
	UpdateReplyCommentContent(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID, replyID primitive.ObjectID, newContent string) error
}

type PostService struct {
	repo PostRepository
}

func NewPostService(repo PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(ctx context.Context, title, content, author string) (model.Post, error) {
	if title == "" {
		return model.Post{}, errors.New("title is required")
	}
	if content == "" {
		return model.Post{}, errors.New("content is required")
	}
	p, err := s.repo.CreatePost(ctx, title, content, author)
	if err != nil {
		return model.Post{}, err
	}
	return p, nil
}

func (s *PostService) UpdatePost(ctx context.Context, id primitive.ObjectID, title, content string) error {
	err := s.repo.UpdatePost(ctx, id, title, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) DeletePost(ctx context.Context, id primitive.ObjectID) error {
	err := s.repo.DeletePost(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) GetPostByID(ctx context.Context, id primitive.ObjectID) (model.Post, error) {
	p, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		return model.Post{}, err
	}
	return p, nil
}

func (s *PostService) GetAllPosts(ctx context.Context, limit int64, aftercreatedAt *time.Time, afterID *primitive.ObjectID) (posts []model.Post, nextcreatedAt *time.Time, nextID *primitive.ObjectID, err error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, nextcreatedAt, nextID, err = s.repo.ListByCursor(ctx, limit, aftercreatedAt, afterID)
	if err != nil {
		return posts, nextcreatedAt, nextID, err
	}

	return posts, nextcreatedAt, nextID, nil
}

func (s *PostService) AddCommentToPost(ctx context.Context, postID primitive.ObjectID, author string, content string) (model.Comment, error) {
	c, err := s.repo.AddCommentToPost(ctx, postID, author, content)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (s *PostService) DeleteComment(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID) error {
	err := s.repo.DeleteComment(ctx, postID, commentID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) UpdateCommentContent(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID, newContent string) error {
	if newContent == "" {
		return errors.New("newContent is required")
	}
	err := s.repo.UpdateComment(ctx, postID, commentID, newContent)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) ReplyComment(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID, author string, content string) error {
	err := s.repo.ReplyComment(ctx, postID, commentID, author, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) DeleteReplyComment(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID, replyID primitive.ObjectID) error {
	err := s.repo.DeleteReplyComment(ctx, postID, commentID, replyID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) UpdateReplyCommentContent(ctx context.Context, postID primitive.ObjectID, commentID primitive.ObjectID, replyID primitive.ObjectID, newContent string) error {
	if newContent == "" {
		return errors.New("newContent is required")
	}
	err := s.repo.UpdateReplyCommentContent(ctx, postID, commentID, replyID, newContent)
	if err != nil {
		return err
	}
	return nil
}
