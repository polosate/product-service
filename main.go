// product-service/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/micro/go-micro"
	pb "github.com/polosate/product-service/proto/product"
	storageProto "github.com/polosate/storage-service/proto/storage"
)


const (
	defaultHost = "datastore:27017"
)

func main() {
	srv := micro.NewService(
		micro.Name("steaks.product.service"),
	)

	srv.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}

	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())

	productsCollection := client.Database("steaks").Collection("products")

	repository := &MongoRepository{productsCollection}
	storageClient := storageProto.NewStorageServiceClient("steaks.storage.service", srv.Client())
	h := &handler{repository, storageClient}

	// Register handlers
	pb.RegisterProductServiceHandler(srv.Server(), h)

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
