package ports

import (
	"context"

	"backend/internal/consumer-service/core/domain/dto"
)

type IConsumerRepo interface {
	GetStores(ctx context.Context, lat string, lon string, radius string) ([]dto.Store, error)
	GetAllProducts(ctx context.Context) ([]dto.Product, error)
	GetProductInfo(ctx context.Context, productID string) (dto.ProductStoresInfo, error)
	GetStoreProducts(ctx context.Context, storeID string, page string, limit string, search string, sort string) ([]dto.ProductInfo, error)
	GetStoreInfo(ctx context.Context, storeID string) (dto.StoreInfo, error)
	AddCommentToStore(ctx context.Context, storeID string, comment *dto.AddCommentRequest) error
	GetStoreComments(ctx context.Context, storeID string) ([]dto.StoreComment, error)
	UpdateStoreCommentVotes(ctx context.Context, storeID string, commentID string, amount string) error
}
