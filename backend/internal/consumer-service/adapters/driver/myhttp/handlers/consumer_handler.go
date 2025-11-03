package handlers

import (
	"context"
	"net/http"

	ports "backend/internal/consumer-service/core/ports/driver"
	"backend/internal/mylogger"
)

type ConsumerHandler struct {
	mylog   mylogger.Logger
	ctx     context.Context
	service ports.IConsumerService
}

func NewConsumerHandler(ctx context.Context, service ports.IConsumerService, mylog mylogger.Logger) *ConsumerHandler {
	return &ConsumerHandler{
		ctx:     ctx,
		service: service,
		mylog:   mylog,
	}
}

func (h *ConsumerHandler) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
