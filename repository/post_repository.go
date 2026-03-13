package repository

import (
	"context"
	"time"
	"week1-postbackend/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepository struct {
	col *mongo.Collection
}

func NewPostRepository(db *mongo.Database) *PostRepository {
	return &PostRepository{
		col: db.Collection("posts"),
	}
}

func (r *PostRepository) CreatePost(ctx context.Context, title, content, author string) (model.Post, error) {
	now := time.Now()
	p := model.Post{
		Title:   title,
		Body:    content,
		Author:  author,
		Created: now,
		Updated: now,
	}
	_, err := r.col.InsertOne(ctx, p)
	if err != nil {
		return model.Post{}, err
	}
	return p, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, id primitive.ObjectID, title, content string) error {
	now := time.Now()
	opts := bson.M{
		"$set": bson.M{
			"title":   title,
			"body":    content,
			"updated": now,
		},
	}
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, opts)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) DeletePost(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) GetPostByID(ctx context.Context, id primitive.ObjectID) (model.Post, error) {
	var p model.Post
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err != nil {
		return model.Post{}, err
	}
	return p, nil
}

func (r *PostRepository) ListByCursor(
	ctx context.Context,
	limit int64,
	afterCreated *time.Time,
	afterID *primitive.ObjectID,
) (posts []model.Post, nextCreatedAt *time.Time, nextID *primitive.ObjectID, err error) {

	filter := bson.M{}
	if afterCreated != nil && afterID != nil {
		filter = bson.M{
			"$or": bson.A{
				bson.M{"createdAt": bson.M{"$lt": *afterCreated}},
				bson.M{
					"createdAt": *afterCreated,
					"_id":       bson.M{"$lt": *afterID},
				},
			},
		}
	}

	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.D{
			{Key: "createdAt", Value: -1},
			{Key: "_id", Value: -1},
		})

	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, nil, err
	}
	defer cur.Close(ctx)

	if err = cur.All(ctx, &posts); err != nil {
		return nil, nil, nil, err
	}

	if len(posts) == 0 {
		return posts, nil, nil, nil
	}
	last := posts[len(posts)-1]
	nc := last.Created
	nid := last.ID

	return posts, &nc, &nid, nil
}

func (r *PostRepository) AddCommentToPost(
	ctx context.Context,
	postID primitive.ObjectID,
	author string,
	content string) (model.Comment, error) {

	now := time.Now()

	c := model.Comment{
		ID:       primitive.NewObjectID(),
		Author:   author,
		Content:  content,
		Comments: []model.Comment{},
		Created:  now,
		Updated:  now,
	}

	filter := bson.M{"_id": postID}
	update := bson.M{"$push": bson.M{"comments": c}}

	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return model.Comment{}, err
	}
	if res.MatchedCount == 0 {
		return c, mongo.ErrNoDocuments
	}

	return c, nil
}

func (r *PostRepository) DeleteComment(
	ctx context.Context,
	postID primitive.ObjectID,
	commentID primitive.ObjectID,
) error {
	filter := bson.M{"_id": postID}
	update := bson.M{"$pull": bson.M{"comments": commentID}}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *PostRepository) UpdateComment(
	ctx context.Context,
	postID primitive.ObjectID,
	commentID primitive.ObjectID,
	newContent string,
) error {
	filter := bson.M{"_id": postID, "comments._id": commentID}
	update := bson.M{
		"$push": bson.M{
			"comments.$.content":   newContent,
			"comments.$.updatedAt": time.Now().UTC()},
	}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *PostRepository) ReplyComment(
	ctx context.Context,
	postID primitive.ObjectID,
	commentID primitive.ObjectID,
	author string,
	content string,
) error {
	now := time.Now()
	reply := model.Comment{
		ID:       primitive.NewObjectID(),
		Author:   author,
		Content:  content,
		Comments: []model.Comment{},
		Created:  now,
		Updated:  now,
	}
	filter := bson.M{"_id": postID}
	update := bson.M{
		"$push": bson.M{
			"comments.$[c].comments": reply,
		},
	}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"c._id": commentID},
		},
	})
	res, err := r.col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	if res.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *PostRepository) DeleteReplyComment(
	ctx context.Context,
	postID primitive.ObjectID,
	commentID primitive.ObjectID,
	replyID primitive.ObjectID,
) error {

	filter := bson.M{"_id": postID}

	update := bson.M{
		"$pull": bson.M{
			"comments.$[p].comments": bson.M{"_id": replyID},
		},
	}

	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"p._id": commentID},
		},
	})

	res, err := r.col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 || res.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *PostRepository) UpdateReplyCommentContent(
	ctx context.Context,
	postID primitive.ObjectID,
	commentID primitive.ObjectID,
	replyID primitive.ObjectID,
	newContent string,
) error {

	now := time.Now().UTC()

	filter := bson.M{"_id": postID}

	update := bson.M{
		"$set": bson.M{
			"comments.$[p].comments.$[c].content":   newContent,
			"comments.$[p].comments.$[c].updatedAt": now,
		},
	}

	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"p._id": commentID},
			bson.M{"c._id": replyID},
		},
	})

	res, err := r.col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 || res.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
