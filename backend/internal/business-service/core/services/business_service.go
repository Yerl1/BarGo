package services

import (
	"context"

	"backend/internal/business-service/core/domain/dto"
	"backend/internal/mylogger"

	ports "backend/internal/business-service/core/ports/driven"
)

type BusinessService struct {
	mylog        mylogger.Logger
	BusinessRepo ports.IBusinessRepo
	ctx          context.Context
	jwtSecret    string
}

func NewBusinessService(ctx context.Context, BusinessRepo ports.IBusinessRepo, jwtSecret string, mylog mylogger.Logger) *BusinessService {
	return &BusinessService{
		ctx:          ctx,
		BusinessRepo: BusinessRepo,
		jwtSecret:    jwtSecret,
		mylog:        mylog,
	}
}

func (s *BusinessService) GetStores(ctx context.Context, userId string) ([]dto.Store, error) {
	return s.BusinessRepo.GetStores(ctx, userId)
}

func (s *BusinessService) AddStore(ctx context.Context, storeId string, store dto.Store) error {
	return s.BusinessRepo.AddStore(ctx, storeId, store)
}

func (s *BusinessService) UpdateStore(ctx context.Context, store dto.Store) error {
	return s.BusinessRepo.UpdateStore(ctx, store)
}

func (s *BusinessService) DeleteStore(ctx context.Context, storeId string) error {
	return s.BusinessRepo.DeleteStore(ctx, storeId)
}

func (s *BusinessService) GetProducts(ctx context.Context, storeId string) ([]dto.Product, error) {
	return s.BusinessRepo.GetProducts(ctx, storeId)
}

func (s *BusinessService) AddProduct(ctx context.Context, storeId string, product dto.Product) error {
	return s.BusinessRepo.AddProduct(ctx, storeId, product)
}

func (s *BusinessService) UpdateProduct(ctx context.Context, product dto.Product) error {
	return s.BusinessRepo.UpdateProduct(ctx, product)
}

func (s *BusinessService) DeleteProduct(ctx context.Context, productId string) error {
	return s.BusinessRepo.DeleteProduct(ctx, productId)
}
