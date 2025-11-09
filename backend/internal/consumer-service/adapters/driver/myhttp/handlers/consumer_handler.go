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

func (h *ConsumerHandler) GetStores() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation for getting stores based on lat, lon, and radius
		lat := r.URL.Query().Get("lat")
		lon := r.URL.Query().Get("lon")
		radius := r.URL.Query().Get("radius")

		// Call the service layer to get stores
		stores, err := h.service.GetStores(h.ctx, lat, lon, radius)
		if err != nil {
			http.Error(w, "Failed to fetch stores", http.StatusInternalServerError)
			return
		}

		h.mylog.Debug().
			Str("lat", lat).
			Str("lon", lon).
			Str("radius", radius).
			Int("store_count", len(stores)).
			Msg("Fetched stores successfully")

		// Serialize and write the response
		JsonResponse(w, stores, http.StatusOK)
	}
}

func (h *ConsumerHandler) GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Call the service layer to get stores
		products, err := h.service.GetAllProducts(h.ctx)
		if err != nil {
			http.Error(w, "Failed to fetch stores", http.StatusInternalServerError)
			return
		}

		h.mylog.Debug().
			Int("products_count", len(products)).
			Msg("Fetched stores successfully")

		// Serialize and write the response
		JsonResponse(w, products, http.StatusOK)
	}
}

func (h *ConsumerHandler) GetProductInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productId := r.URL.Query().Get("product_id")
		// Call the service layer to get stores
		productInfo, err := h.service.GetProductInfo(h.ctx, productId)
		if err != nil {
			http.Error(w, "Failed to fetch stores", http.StatusInternalServerError)
			return
		}

		h.mylog.Debug().
			Msg("Fetched productInfo successfully")

		// Serialize and write the response
		JsonResponse(w, productInfo, http.StatusOK)
	}
}

func (h *ConsumerHandler) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
