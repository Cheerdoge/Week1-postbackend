package mongodb

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoClient() (*mongo.Client, error) {
	uri := os.Getenv("MONGO_URI")

	//配置
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	//连接建立客户端
	client, err := mongo.Connect(opts)
	if err != nil {
		fmt.Println("Error connecting to MongoDB,")
		return nil, err
	}
	//验证连接
	var result bson.M
	if err := client.Database("admin").
		RunCommand(context.TODO(), bson.D{{"ping", 1}}).
		Decode(&result); err != nil {
		return nil, err
	}
	return client, nil
}
