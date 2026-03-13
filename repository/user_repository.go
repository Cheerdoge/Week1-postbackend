package repository

import (
	"context"
	"time"
	"week1-postbackend/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		col: db.Collection("user"),
	}
}

func (r UserRepository) CreateUser(ctx context.Context, username, email, password string) (model.User, error) {
	now := time.Now()
	u := model.User{
		Username: username,
		Email:    email,
		Password: password,
		Created:  now,
	}
	_, err := r.col.InsertOne(ctx, u)
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (r UserRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (model.User, error) {
	var u model.User
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (r UserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var u model.User
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}
