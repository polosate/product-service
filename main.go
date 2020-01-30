// steaks/main.go
package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"log"
	"sync"

	// Import the generated protobuf code
	pb "github.com/polosate/steaks/proto/product"
	storageProto "github.com/polosate/steaks-service-storage/proto/storage"
)

type repository interface {
	Create(*pb.Product) (*pb.Product, error)
	GetAll() []*pb.Product
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	mu sync.RWMutex
	products []*pb.Product
}

// Create a new products
func (repo *Repository) Create(product *pb.Product) (*pb.Product, error) {
	repo.mu.Lock()
	updated := append(repo.products, product)
	repo.products = updated
	repo.mu.Unlock()
	return product, nil
}

func (repo *Repository) GetAll() []*pb.Product {
	return repo.products
}

type service struct {
	repo repository
	storageClient storageProto.StorageServiceClient
}

// CreateProduct - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *service) CreateProduct(ctx context.Context, req *pb.Product, res *pb.Response) error {

	storageResponse, err := s.storageClient.FindAvailable(context.Background(), &storageProto.Specification{
		Capacity: req.Amount,
	})
	if err != nil {
		return err
	}
	log.Printf("Found storage: %s \n", storageResponse.Storage.Name)

	req.StorageId = storageResponse.Storage.Id

	// Save our consignment
	product, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Product = product
	return nil
}

func (s *service) GetProducts(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	products := s.repo.GetAll()
	res.Products = products
	return nil
}

func main() {

	repo := &Repository{}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("steaks.product.service"),
	)

	// Init will parse the command line flags.
	srv.Init()

	storageClient := storageProto.NewStorageServiceClient("steaks.service.storages", srv.Client())

	// Register handler
	pb.RegisterProductServiceHandler(srv.Server(), &service{repo, storageClient})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
