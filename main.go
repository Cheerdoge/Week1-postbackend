package main

import (
	"context"
	"log"
	"os"
	"week1-postbackend/handler"
	"week1-postbackend/infra/mongodb"
	"week1-postbackend/repository"
	"week1-postbackend/router"
	"week1-postbackend/service"
)

func main() {
	ctx := context.Background()

	client, err := mongodb.NewMongoClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "week1"
	}
	db := client.Database(dbName)

	// repo
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)

	// service
	userSvc := service.NewUserService(userRepo)
	postSvc := service.NewPostService(postRepo)

	// handler
	authHandler := handler.NewAuthHandler(userSvc)
	postHandler := handler.NewPostHandler(postSvc)

	// router
	r := router.NewRouter(router.Handlers{
		Auth: authHandler,
		Post: postHandler,
	})
	
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
