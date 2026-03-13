package service

import (
	"context"
	"week1-postbackend/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username, email, password string) (model.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type UserService struct {
	repo UserRepository
}
