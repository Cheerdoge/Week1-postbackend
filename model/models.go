package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`

	Created time.Time  `bson:"createdAt"`
	Deleted *time.Time `bson:"deletedAt"`
}

type Comment struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Author   string             `bson:"author"`
	Content  string             `bson:"content"`
	Comments []Comment          `bson:"comments"`

	Created time.Time `bson:"createdAt"`
	Updated time.Time `bson:"updatedAt"`
}

type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Title    string             `bson:"title"`
	Body     string             `bson:"body"`
	Author   string             `bson:"author"`
	Comments []Comment          `bson:"comments"`

	Created time.Time `bson:"createdAt"`
	Updated time.Time `bson:"updatedAt"`
}
