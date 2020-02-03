// product-service/handler.go
package main

import (
	"context"
	"github.com/pkg/errors"
	pb "github.com/polosate/product-service/proto/product"
	storageProto "github.com/polosate/storage-service/proto/storage"
)

type handler struct {
	repository
	storageClient storageProto.StorageServiceClient
}

func (s *handler) CreateProduct(ctx context.Context, req *pb.Product, res *pb.Response) error {

	storageResponse, err := s.storageClient.FindAvailable(ctx, &storageProto.Specification{
		Capacity: req.Amount,
	})
	if storageResponse == nil {
		return errors.New("error fetching storage, returned nil")
	}

	if err != nil {
		return err
	}

	req.StorageId = storageResponse.Storage.Id

	if err = s.repository.Create(ctx, MarshalProduct(req)); err != nil {
		return err
	}

	res.Created = true
	res.Product = req
	return nil
}

func (s *handler) GetProducts(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	products, err := s.repository.GetAll(ctx)
	if err != nil {
		return err
	}
	res.Products = UnmarshalProductCollection(products)
	return nil
}

