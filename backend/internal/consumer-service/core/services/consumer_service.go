package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"backend/internal/consumer-service/core/domain/dto"
	"backend/internal/mylogger"

	ports "backend/internal/consumer-service/core/ports/driven"
)

type ConsumerService struct {
	mylog        mylogger.Logger
	ConsumerRepo ports.IConsumerRepo
	ctx          context.Context
	jwtSecret    string

	fastAPIURL string
	client     *http.Client
}

func NewConsumerService(ctx context.Context, consumerRepo ports.IConsumerRepo, jwtSecret string, mylog mylogger.Logger) *ConsumerService {
	return &ConsumerService{
		ctx:          ctx,
		ConsumerRepo: consumerRepo,
		jwtSecret:    jwtSecret,
		mylog:        mylog,
		fastAPIURL:   "http://ml-fastapi-service:3005/predict",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *ConsumerService) GetStores(ctx context.Context, lat string, lon string, radius string) ([]dto.Store, error) {
	return s.ConsumerRepo.GetStores(ctx, lat, lon, radius)
}

func (s *ConsumerService) GetAllProducts(ctx context.Context) ([]dto.Product, error) {
	return s.ConsumerRepo.GetAllProducts(ctx)
}

func (s *ConsumerService) GetProductInfo(ctx context.Context, productID string) (dto.ProductStoresInfo, error) {
	return s.ConsumerRepo.GetProductInfo(ctx, productID)
}

func (s *ConsumerService) GetStoreProducts(ctx context.Context, storeID string, page string, limit string, search string, sort string) ([]dto.ProductInfo, error) {
	return s.ConsumerRepo.GetStoreProducts(ctx, storeID, page, limit, search, sort)
}

func (s *ConsumerService) GetStoreInfo(ctx context.Context, storeID string) (dto.StoreInfo, error) {
	return s.ConsumerRepo.GetStoreInfo(ctx, storeID)
}

func (s *ConsumerService) AddCommentToStore(
	ctx context.Context,
	storeID string,
	comment *dto.AddCommentRequest,
) error {
	log := s.mylog.With().Str("service", "AddCommentToStore").Logger()

	if comment == nil || comment.Content == "" {
		log.Error().Msg("Comment content is empty")
		return fmt.Errorf("comment is empty")
	}

	if len(comment.Content) > 500 {
		log.Error().Msg("Comment exceeds maximum length")
		return fmt.Errorf("comment exceeds maximum length of 500 characters")
	}

	if comment.UserID == "" {
		log.Error().Msg("User ID is missing in comment")
		return fmt.Errorf("user ID is required")
	}

	if comment.Rating < 1 || comment.Rating > 5 {
		log.Error().Msg("Invalid rating value")
		return fmt.Errorf("rating must be between 1 and 5")
	}

	// Prepare request body for ML service
	toxicityReq := dto.PredictionRequest{
		Comment: comment.Content,
	}

	jsonBody, err := json.Marshal(toxicityReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal toxicity request")
		return fmt.Errorf("failed to marshal toxicity request: %w", err)
	}

	// Correct ML FastAPI URL
	mlURL := "http://bargo-ml-api:3005/predict"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, mlURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error().Err(err).Msg("Creating toxicity request failed")
		return fmt.Errorf("creating toxicity request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Toxicity service request failed")
		return fmt.Errorf("toxicity service request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Str("status", resp.Status).Msg("Toxicity service returned non-200 status")
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("toxicity service error: status=%d body=%s", resp.StatusCode, string(body))
	}

	// Correct struct with "toxic"
	var prediction dto.PredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		log.Error().Err(err).Msg("Failed to decode toxicity response")
		return fmt.Errorf("failed to decode toxicity response: %w", err)
	}

	if prediction.IsToxic {
		log.Error().Msg("Comment detected as toxic, cannot add")
		return fmt.Errorf("comment detected as toxic, cannot add")
	}

	// Add comment to DB
	if err := s.ConsumerRepo.AddCommentToStore(ctx, storeID, comment); err != nil {
		log.Error().Err(err).Msg("Failed to add comment to store")
		return fmt.Errorf("failed to add comment to store: %w", err)
	}

	log.Info().Msg("Comment added to store successfully")
	return nil
}

func (s *ConsumerService) GetStoreComments(ctx context.Context, storeID string) ([]dto.StoreComment, error) {
	return s.ConsumerRepo.GetStoreComments(ctx, storeID)
}

func (s *ConsumerService) UpdateStoreCommentVotes(ctx context.Context, storeID string, commentID string, amount string) (dto.StoreInfo, error) {
	return s.ConsumerRepo.UpdateStoreCommentVotes(ctx, storeID, commentID, amount)
}
