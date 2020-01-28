// steaks/main.go
package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	// Import the generated protobuf code
	pb "github.com/polosate/steaks/proto/product"
)

type repository interface {
	Create(*pb.Product) (*pb.Product, error)
	GetAll() []*pb.Product
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	products []*pb.Product
}

// Create a new products
func (repo *Repository) Create(product *pb.Product) (*pb.Product, error) {
	updated := append(repo.products, product)
	repo.products = updated
	return product, nil
}

func (repo *Repository) GetAll() []*pb.Product {
	return repo.products
}

type service struct {
	repo repository
}

// CreateProduct - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *service) CreateProduct(ctx context.Context, req *pb.Product, res *pb.Response) error {

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

	// Register handler
	pb.RegisterProductServiceHandler(srv.Server(), &service{repo})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
