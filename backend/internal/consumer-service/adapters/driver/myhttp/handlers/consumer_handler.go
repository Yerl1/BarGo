package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"backend/internal/consumer-service/core/domain/dto"
	ports "backend/internal/consumer-service/core/ports/driver"
	"backend/internal/mylogger"

	"github.com/rs/zerolog/log"
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
			JsonResponse(w, "Failed to fetch all product stores", http.StatusInternalServerError)
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
			JsonResponse(w, "Failed to fetch all products", http.StatusInternalServerError)
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
		productStoresInfo, err := h.service.GetProductInfo(h.ctx, productId)
		if err != nil {
			JsonResponse(w, "Failed to fetch product info", http.StatusInternalServerError)
			return
		}

		h.mylog.Debug().
			Msg("Fetched productStoresInfo successfully")

		// Serialize and write the response
		JsonResponse(w, productStoresInfo, http.StatusOK)
	}
}

func (h *ConsumerHandler) GetStoreProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storeId := r.PathValue("store_id")

		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")
		search := r.URL.Query().Get("search")
		sort := r.URL.Query().Get("sort")

		// Call the service layer to get stores
		products, err := h.service.GetStoreProducts(h.ctx, storeId, page, limit, search, sort)
		if err != nil {
			JsonResponse(w, "Failed to fetch store products", http.StatusInternalServerError)
			return
		}

		h.mylog.Debug().
			Int("products_count", len(products)).
			Msg("Fetched stores successfully")

		// Serialize and write the response
		JsonResponse(w, products, http.StatusOK)
	}
}

func (h *ConsumerHandler) GetStoreInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storeId := r.PathValue("store_id")

		// Call the service layer to get stores
		storeInfo, err := h.service.GetStoreInfo(h.ctx, storeId)
		if err != nil {
			JsonResponse(w, "Failed to fetch store info", http.StatusInternalServerError)
			return
		}

		h.mylog.Debug().
			Str("store_name", storeInfo.Name).
			Msg("Store info fetched successfully")

		// Serialize and write the response
		JsonResponse(w, storeInfo, http.StatusOK)
	}
}

func (h *ConsumerHandler) AddCommentToStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var comment dto.AddCommentRequest

		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			log.Error().Err(err).Msg("Failed to decode request body")
			JsonResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		storeId := r.PathValue("store_id")
		// Call the service layer to get stores
		err := h.service.AddCommentToStore(h.ctx, storeId, &comment)
		if err != nil {
			JsonResponse(w, "Failed to add comment to store", http.StatusInternalServerError)
			return
		}

		h.mylog.Debug().Msg("Comment added successfully!")

		// Serialize and write the response
		JsonResponse(w, "Comment added successfully!", http.StatusOK)
	}
}

func (h *ConsumerHandler) GetStoreComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storeId := r.PathValue("store_id")
		// Call the service layer to get stores
		storeComments, err := h.service.GetStoreComments(h.ctx, storeId)
		if err != nil {
			JsonResponse(w, "Failed to fetch store comments", http.StatusInternalServerError)
			return
		}

		h.mylog.Debug().Msg("Store comments fetched successfully")

		// Serialize and write the response
		JsonResponse(w, storeComments, http.StatusOK)
	}
}

func (h *ConsumerHandler) UpdateStoreCommentVotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storeId := r.PathValue("store_id")
		commentId := r.PathValue("comment_id")

		amount := r.URL.Query().Get("amount") // e.g., "1" or "-1"

		// Call the service layer to get stores
		err := h.service.UpdateStoreCommentVotes(h.ctx, storeId, commentId, amount)
		if err != nil {
			JsonResponse(w, "Failed to update store comment votes", http.StatusInternalServerError)
			return
		}

		h.mylog.Debug().
			Msg("Store votes updated successfully")

		// Serialize and write the response
		JsonResponse(w, "Store votes updated successfully", http.StatusOK)
	}
}

func (h *ConsumerHandler) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
