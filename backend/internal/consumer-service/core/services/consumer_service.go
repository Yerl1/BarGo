package services

import (
	"context"

	"backend/internal/consumer-service/core/domain/dto"
	"backend/internal/mylogger"

	ports "backend/internal/consumer-service/core/ports/driven"
)

type ConsumerService struct {
	mylog        mylogger.Logger
	ConsumerRepo ports.IConsumerRepo
	ctx          context.Context
	jwtSecret    string
}

func NewConsumerService(ctx context.Context, consumerRepo ports.IConsumerRepo, jwtSecret string, mylog mylogger.Logger) *ConsumerService {
	return &ConsumerService{
		ctx:          ctx,
		ConsumerRepo: consumerRepo,
		jwtSecret:    jwtSecret,
		mylog:        mylog,
	}
}

func (s *ConsumerService) GetStores(ctx context.Context, lat string, lon string, radius string) ([]dto.Store, error) {
	return s.ConsumerRepo.GetStores(ctx, lat, lon, radius)
}

func (s *ConsumerService) GetAllProducts(ctx context.Context) ([]dto.Product, error) {
	return s.ConsumerRepo.GetAllProducts(ctx)
}

func (s *ConsumerService) GetProductInfo(ctx context.Context, productID string) (dto.ProductInfo, error) {
	return s.ConsumerRepo.GetProductInfo(ctx, productID)
}
