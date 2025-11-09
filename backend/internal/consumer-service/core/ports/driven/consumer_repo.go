package ports

import (
	"context"

	"backend/internal/consumer-service/core/domain/dto"
)

type IConsumerRepo interface {
	GetStores(ctx context.Context, lat string, lon string, radius string) ([]dto.Store, error)
	GetAllProducts(ctx context.Context) ([]dto.Product, error)
	GetProductInfo(ctx context.Context, productID string) (dto.ProductInfo, error)
}
