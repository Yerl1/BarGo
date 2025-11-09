package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"backend/internal/business-service/core/domain/dto"
	ports "backend/internal/business-service/core/ports/driver"
	"backend/internal/mylogger"
)

type BusinessHandler struct {
	mylog   mylogger.Logger
	ctx     context.Context
	service ports.IBusinessService
}

func NewBusinessHandler(ctx context.Context, service ports.IBusinessService, mylog mylogger.Logger) *BusinessHandler {
	return &BusinessHandler{
		ctx:     ctx,
		service: service,
		mylog:   mylog,
	}
}

// GetStores — fetch stores belonging to a user
func (h *BusinessHandler) GetStores() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user_id").(string)
		h.mylog.Info().Str("user_id", userId).Msg("GetStores called")

		stores, err := h.service.GetStores(h.ctx, userId)
		if err != nil {
			h.mylog.Error().Err(err).Msg("Failed to fetch stores")
			http.Error(w, "Failed to fetch stores", http.StatusInternalServerError)
			return
		}

		h.mylog.Info().Msg("Fetched stores successfully")
		JsonResponse(w, stores, http.StatusOK)
	}
}

// AddStore — creates a new store for the user
func (h *BusinessHandler) AddStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user_id").(string)

		var store dto.Store
		if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		h.mylog.Info().Str("user_id", userId).Str("store_name", *store.Name).Msg("AddStore called")

		if err := h.service.AddStore(h.ctx, userId, store); err != nil {
			h.mylog.Error().Err(err).Msg("Failed to add store")
			http.Error(w, "Failed to add store", http.StatusInternalServerError)
			return
		}

		h.mylog.Info().Msg("Store added successfully")
		JsonResponse(w, map[string]string{"message": "Store added successfully"}, http.StatusCreated)
	}
}

// UpdateStore — modifies existing store
func (h *BusinessHandler) UpdateStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var store dto.Store
		if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		h.mylog.Info().Str("store_id", store.StoreID).Msg("UpdateStore called")

		if err := h.service.UpdateStore(h.ctx, store); err != nil {
			h.mylog.Error().Err(err).Msg("Failed to update store")
			http.Error(w, "Failed to update store", http.StatusInternalServerError)
			return
		}

		JsonResponse(w, map[string]string{"message": "Store updated successfully"}, http.StatusOK)
	}
}

// DeleteStore — removes a store by ID
func (h *BusinessHandler) DeleteStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storeId := r.URL.Query().Get("store_id")
		if storeId == "" {
			http.Error(w, "store_id is required", http.StatusBadRequest)
			return
		}

		h.mylog.Info().Str("store_id", storeId).Msg("DeleteStore called")

		if err := h.service.DeleteStore(h.ctx, storeId); err != nil {
			h.mylog.Error().Err(err).Msg("Failed to delete store")
			http.Error(w, "Failed to delete store", http.StatusInternalServerError)
			return
		}

		JsonResponse(w, map[string]string{"message": "Store deleted successfully"}, http.StatusOK)
	}
}

// GetProducts — list all products for a store
func (h *BusinessHandler) GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storeId := r.URL.Query().Get("store_id")
		if storeId == "" {
			http.Error(w, "store_id is required", http.StatusBadRequest)
			return
		}

		h.mylog.Info().Str("store_id", storeId).Msg("GetProducts called")

		products, err := h.service.GetProducts(h.ctx, storeId)
		if err != nil {
			h.mylog.Error().Err(err).Msg("Failed to fetch products")
			http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
			return
		}

		JsonResponse(w, products, http.StatusOK)
	}
}

// AddProduct — create new product
func (h *BusinessHandler) AddProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product dto.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		h.mylog.Info().Str("product_name", product.Name).Msg("AddProduct called")

		storeId := r.URL.Query().Get("store_id")

		if err := h.service.AddProduct(h.ctx, storeId, product); err != nil {
			h.mylog.Error().Err(err).Msg("Failed to add product")
			http.Error(w, "Failed to add product", http.StatusInternalServerError)
			return
		}

		JsonResponse(w, map[string]string{"message": "Product added successfully"}, http.StatusCreated)
	}
}

// UpdateProduct — edit existing product
func (h *BusinessHandler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product dto.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		h.mylog.Info().Str("product_id", product.ProductID).Msg("UpdateProduct called")

		if err := h.service.UpdateProduct(h.ctx, product); err != nil {
			h.mylog.Error().Err(err).Msg("Failed to update product")
			http.Error(w, "Failed to update product", http.StatusInternalServerError)
			return
		}

		JsonResponse(w, map[string]string{"message": "Product updated successfully"}, http.StatusOK)
	}
}

// DeleteProduct — remove a product by ID
func (h *BusinessHandler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productId := r.URL.Query().Get("product_id")
		if productId == "" {
			http.Error(w, "product_id is required", http.StatusBadRequest)
			return
		}

		h.mylog.Info().Str("product_id", productId).Msg("DeleteProduct called")

		if err := h.service.DeleteProduct(h.ctx, productId); err != nil {
			h.mylog.Error().Err(err).Msg("Failed to delete product")
			http.Error(w, "Failed to delete product", http.StatusInternalServerError)
			return
		}

		JsonResponse(w, map[string]string{"message": "Product deleted successfully"}, http.StatusOK)
	}
}

// HealthHandler — simple readiness probe
func (h *BusinessHandler) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
