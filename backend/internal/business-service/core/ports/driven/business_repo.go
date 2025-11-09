package ports

import (
	"context"

	"backend/internal/business-service/core/domain/dto"
)

type IBusinessRepo interface {
	GetStores(ctx context.Context, userID string) ([]dto.Store, error)
	AddStore(ctx context.Context, userId string, store dto.Store) error
	UpdateStore(ctx context.Context, store dto.Store) error
	DeleteStore(ctx context.Context, storeID string) error

	GetProducts(ctx context.Context, storeID string) ([]dto.Product, error)
	AddProduct(ctx context.Context, storeId string, product dto.Product) error
	UpdateProduct(ctx context.Context, product dto.Product) error
	DeleteProduct(ctx context.Context, productID string) error
}
