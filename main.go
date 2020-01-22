// steaks/main.go
package main

import (
	"context"
	"log"
	"net"
	"sync"

	// Import the generated protobuf code
	pb "github.com/polosate/steaks/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Product) (*pb.Product, error)
	GetAll() []*pb.Product
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	mu       sync.RWMutex
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

// GetAll returns all products
func (repo *Repository) GetAll() []*pb.Product {
	return repo.products
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo repository
}

// CreateProduct - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *service) CreateProduct(ctx context.Context, req *pb.Product) (*pb.Response, error) {

	// Save our consignment
	product, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	return &pb.Response{Created: true, Product: product}, nil
}

func (s *service) GetProducts(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	products := s.repo.GetAll()
	return &pb.Response{Products: products}, nil
}

func main() {

	repo := &Repository{}

	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterProductServiceServer(s, &service{repo})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
