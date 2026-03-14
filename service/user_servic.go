package service

import (
	"context"
	"errors"
	"time"
	"week1-postbackend/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username, email, password string) (model.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredential  = errors.New("invalid email or password")
)

func (s *UserService) Register(ctx context.Context, username, email, password string) (model.User, error) {
	if username == "" || email == "" || password == "" {
		return model.User{}, errors.New("username/email/password is required")
	}

	// email 已存在
	_, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return model.User{}, ErrEmailAlreadyExists
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	u, err := s.repo.CreateUser(ctx, username, email, string(hash))
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (model.User, error) {
	if email == "" || password == "" {
		return model.User{}, errors.New("email/password is required")
	}

	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, ErrInvalidCredential
		}
		return model.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return model.User{}, ErrInvalidCredential
	}

	// 可以顺手更新 lastLoginAt（你 model 里暂时没有，就不写）
	_ = time.Now().UTC()

	return u, nil
}
