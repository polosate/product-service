// product-service/repository.go
package main

import (
	"context"
	pb "github.com/polosate/product-service/proto/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Product struct {
	Id int32 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Amount int32 `json:"amount"`
	StorageId int32 `json:"storage_id"`
}

func UnmarshalProductCollection(products []*Product) []*pb.Product {
	collection := make([]*pb.Product, 0)
	for _, product := range products {
		collection = append(collection, UnmarshalProduct(product))
	}
	return collection
}

func MarshalProduct(product *pb.Product) *Product {
	return &Product{
		Id: product.Id,
		Name: product.Name,
		Description: product.Description,
		Amount: product.Amount,
		StorageId: product.StorageId,
	}
}

func UnmarshalProduct(product *Product) *pb.Product {
	return &pb.Product{
		Id: product.Id,
		Name: product.Name,
		Description: product.Description,
		Amount: product.Amount,
		StorageId: product.StorageId,
	}
}

type repository interface {
	Create(ctx context.Context, product *Product) error
	GetAll(ctx context.Context) ([]*Product, error)
}

// MongoRepository implementation
type MongoRepository struct {
	collection *mongo.Collection
}

// Create -
func (repository *MongoRepository) Create(ctx context.Context, product *Product) error {
	_, err := repository.collection.InsertOne(ctx, product)
	return err
}

// GetAll -
func (repository *MongoRepository) GetAll(ctx context.Context) ([]*Product, error) {
	cur, err := repository.collection.Find(ctx, bson.D{}, nil)
	var products []*Product
	for cur.Next(ctx) {
		var product *Product
		if err := cur.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, err
}

