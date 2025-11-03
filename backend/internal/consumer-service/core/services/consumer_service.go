package services

import (
	"context"

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
