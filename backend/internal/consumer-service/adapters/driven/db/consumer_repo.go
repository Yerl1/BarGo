package db

import (
	"context"
	"fmt"
	"strconv"

	"backend/internal/consumer-service/core/domain/dto"
	"backend/internal/mylogger"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type ConsumerRepo struct {
	ctx   context.Context
	mylog mylogger.Logger
	DB    *DB
}

func NewConsumerRepo(ctx context.Context, db *DB, mylog mylogger.Logger) *ConsumerRepo {
	return &ConsumerRepo{
		ctx:   ctx,
		mylog: mylog,
		DB:    db,
	}
}

// GetStores returns all stores within the specified radius (in meters)
// from the given latitude and longitude.
func (r *ConsumerRepo) GetStores(ctx context.Context, lat, lon, radius string) ([]dto.Store, error) {
	log := log.With().Str("method", "GetStores").Logger()
	log.Debug().
		Str("lat", lat).
		Str("lon", lon).
		Str("radius", radius).
		Msg("fetching nearby stores")

	query := `
		SELECT
				s.store_id,
				s.name,
				s.address,
				c.latitude,
				c.longitude
		FROM stores s
		JOIN coords c ON s.coord = c.coord_id
		WHERE ST_DWithin(
				ST_SetSRID(ST_MakePoint(c.longitude, c.latitude), 4326)::geography,
				ST_SetSRID(ST_MakePoint($1::double precision, $2::double precision), 4326)::geography,
				$3::double precision
		);

	`

	latF, err := strconv.ParseFloat(lat, 64)
	lonF, err := strconv.ParseFloat(lon, 64)
	radiusF, err := strconv.ParseFloat(radius, 64)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.conn.Query(ctx, query, lonF, latF, radiusF)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", "GetStores").
			Str("lat", lat).
			Str("lon", lon).
			Str("radius", radius).
			Msg("failed to query nearby stores")
		return nil, err
	}
	defer rows.Close()

	var stores []dto.Store
	for rows.Next() {
		var s dto.Store
		if err := rows.Scan(
			&s.StoreID,
			&s.Name,
			&s.Address,
			&s.Latitude,
			&s.Longitude,
		); err != nil {
			log.Error().Err(err).Msg("failed to scan store row")
			return nil, err
		}
		stores = append(stores, s)
	}

	if err = rows.Err(); err != nil {
		log.Error().
			Err(err).
			Str("query", "GetStores").
			Msg("row iteration error")
		return nil, err
	}

	log.Debug().
		Int("count", len(stores)).
		Str("lat", lat).
		Str("lon", lon).
		Str("radius", radius).
		Msg("successfully fetched nearby stores")

	return stores, nil
}

func (r *ConsumerRepo) GetAllProducts(ctx context.Context) ([]dto.Product, error) {
	log := log.With().Str("method", "GetAllProducts").Logger()
	log.Debug().Msg("fetching all products")

	query := `
		SELECT
			product_id,
			name,
			description,
			photo,
			price
		FROM products;
	`

	rows, err := r.DB.conn.Query(ctx, query)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", "GetAllProducts").
			Msg("failed to query products")
		return nil, err
	}
	defer rows.Close()

	var products []dto.Product
	for rows.Next() {
		var p dto.Product
		if err := rows.Scan(
			&p.ProductID,
			&p.Name,
			&p.Description,
			&p.Photo,
			&p.Price,
		); err != nil {
			log.Error().Err(err).Msg("failed to scan product row")
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		log.Error().
			Err(err).
			Str("query", "GetAllProducts").
			Msg("row iteration error")
		return nil, err
	}

	log.Debug().
		Int("count", len(products)).
		Msg("successfully fetched all products")

	return products, nil
}

func (r *ConsumerRepo) GetProductInfo(ctx context.Context, productID string) (dto.ProductStoresInfo, error) {
	log := log.With().Str("method", "GetProductInfo").Logger()
	log.Debug().
		Str("product_id", productID).
		Msg("fetching product info")

	var product dto.Product
	queryProduct := `
		SELECT
			product_id,
			name,
			description,
			photo,
			price
		FROM products
		WHERE product_id = $1;
	`

	err := r.DB.conn.QueryRow(ctx, queryProduct, productID).Scan(
		&product.ProductID,
		&product.Name,
		&product.Description,
		&product.Photo,
		&product.Price,
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", "GetProductInfo - product").
			Str("product_id", productID).
			Msg("failed to query product info")
		return dto.ProductStoresInfo{}, err
	}

	var stores []dto.Store
	queryStores := `
		SELECT
			s.store_id,
			s.name,
			s.address,
			c.latitude,
			c.longitude
		FROM stores s
		JOIN coords c ON s.coord = c.coord_id
		JOIN products p ON p.store_id = s.store_id
		WHERE p.product_id = $1;
	`

	rows, err := r.DB.conn.Query(ctx, queryStores, productID)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", "GetProductInfo - stores").
			Str("product_id", productID).
			Msg("failed to query stores for product")
		return dto.ProductStoresInfo{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var s dto.Store
		if err := rows.Scan(
			&s.StoreID,
			&s.Name,
			&s.Address,
			&s.Latitude,
			&s.Longitude,
		); err != nil {
			log.Error().Err(err).Msg("failed to scan store row for product")
			return dto.ProductStoresInfo{}, err
		}
		stores = append(stores, s)
	}

	if err = rows.Err(); err != nil {
		log.Error().
			Err(err).
			Str("query", "GetProductInfo - stores").
			Msg("row iteration error")
		return dto.ProductStoresInfo{}, err
	}

	log.Debug().
		Str("product_id", productID).
		Int("store_count", len(stores)).
		Msg("successfully fetched product info")

	return dto.ProductStoresInfo{
		Product: product,
		Stores:  stores,
	}, nil
}

func (r *ConsumerRepo) GetStoreInfo(ctx context.Context, storeID string) (dto.StoreInfo, error) {
	query := `
		SELECT 
			s.store_id,
			s.name,
			s.address,
			s.photo,
			s.description,
			s.owner_id,
			c.latitude,
			c.longitude,
			COALESCE(AVG(cm.rating), 0) as average_rating,
			COUNT(cm.comment_id) as review_count
		FROM stores s
		LEFT JOIN coords c ON s.coord = c.coord_id
		LEFT JOIN comments cm ON s.store_id = cm.store_id
		WHERE s.store_id = $1
		GROUP BY 
			s.store_id, s.name, s.address, s.photo, s.description, s.owner_id,
			c.latitude, c.longitude
	`

	var store dto.StoreInfo
	err := r.DB.conn.QueryRow(ctx, query, storeID).Scan(
		&store.StoreID,
		&store.Name,
		&store.Address,
		&store.Photo,
		&store.Description,
		&store.OwnerID,
		&store.Latitude,
		&store.Longitude,
		&store.AverageRating,
		&store.ReviewCount,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return dto.StoreInfo{}, fmt.Errorf("store not found: %s", storeID)
		}
		return dto.StoreInfo{}, fmt.Errorf("failed to get store info: %w", err)
	}

	return store, nil
}

func (r *ConsumerRepo) GetStoreProducts(
	ctx context.Context,
	storeID string,
	page string,
	limit string,
	search string,
	sort string,
) ([]dto.ProductInfo, error) {
	log := r.mylog.Logger.With().
		Str("method", "GetStoreProducts").
		Str("store_id", storeID).
		Str("page", page).
		Str("limit", limit).
		Str("search", search).
		Str("sort", sort).
		Logger()

	log.Info().Msg("Fetching store products")

	// ---------------------
	// Parse pagination
	// ---------------------
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		log.Warn().Err(err).Msg("Invalid page param, defaulting to 1")
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 100 {
		log.Warn().Err(err).Msg("Invalid limit param, defaulting to 12")
		limitInt = 12
	}

	offset := (pageInt - 1) * limitInt

	// ---------------------
	// Base query
	// ---------------------
	baseQuery := `
		SELECT 
			product_id,
			store_id,
			name,
			description,
			photo,
			price,
			created_at,
			updated_at,
			(price > 0) as in_stock
		FROM products
		WHERE store_id = $1
	`

	args := []interface{}{storeID}
	argPos := 2 // Next placeholder index

	// ---------------------
	// Search filter
	// ---------------------
	if search != "" {
		log.Debug().Msg("Applying search filter to query")

		baseQuery += fmt.Sprintf(`
			AND (
				name ILIKE $%d
				OR description ILIKE $%d
			)
		`, argPos, argPos)

		args = append(args, "%"+search+"%")
		argPos++
	}

	// ---------------------
	// Sorting
	// ---------------------
	orderBy := "created_at DESC" // Default sorting

	switch sort {
	case "price_asc":
		orderBy = "price ASC"

	case "price_desc":
		orderBy = "price DESC"

	case "relevance":
		if search != "" {
			log.Debug().Msg("Applying relevance sorting")

			// Fix: relevance must use the correct arg index
			orderBy = fmt.Sprintf(`
				CASE 
					WHEN name ILIKE $%d THEN 1
					WHEN description ILIKE $%d THEN 2
					ELSE 3
				END, created_at DESC
			`, argPos-1, argPos-1)
		}
	}

	// ---------------------
	// Final query assembly
	// ---------------------
	query := fmt.Sprintf(`
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, baseQuery, orderBy, argPos, argPos+1)

	args = append(args, limitInt, offset)

	log.Debug().
		Str("query", query).
		Int("args_count", len(args)).
		Interface("args", args).
		Msg("Executing SQL query")

	// ---------------------
	// Execute query
	// ---------------------
	rows, err := r.DB.conn.Query(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Query execution failed")
		return nil, fmt.Errorf("failed to query store products: %w", err)
	}
	defer rows.Close()

	// ---------------------
	// Process rows
	// ---------------------
	var products []dto.ProductInfo

	for rows.Next() {
		var product dto.ProductInfo

		if err := rows.Scan(
			&product.ProductID,
			&product.StoreID,
			&product.Name,
			&product.Description,
			&product.Photo,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.InStock,
		); err != nil {
			log.Error().Err(err).Msg("Failed to scan row")
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		products = append(products, product)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Row iteration error")
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	log.Info().
		Int("count", len(products)).
		Msg("Products fetched successfully")

	return products, nil
}

func (r *ConsumerRepo) AddCommentToStore(ctx context.Context, storeID string, comment *dto.AddCommentRequest) error {
	query := `
		INSERT INTO comments (store_id, user_id, content, rating, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW());
	`

	_, err := r.DB.conn.Exec(ctx, query,
		storeID,
		comment.UserID,
		comment.Content,
		comment.Rating,
	)
	if err != nil {
		return fmt.Errorf("failed to add comment to store: %w", err)
	}

	return nil
}

func (r *ConsumerRepo) GetStoreComments(ctx context.Context, storeID string) ([]dto.StoreComment, error) {
	query := `
		SELECT
			comment_id,
			user_id,
			content,
			rating,
			helpful_votes
		FROM comments
		WHERE store_id = $1;
	`

	rows, err := r.DB.conn.Query(ctx, query, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query store comments: %w", err)
	}
	defer rows.Close()

	var comments []dto.StoreComment
	for rows.Next() {
		var comment dto.StoreComment
		if err := rows.Scan(
			&comment.CommentID,
			&comment.UserID,
			&comment.Content,
			&comment.Rating,
			&comment.HelpfulVotes,
		); err != nil {
			return nil, fmt.Errorf("failed to scan store comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating store comments: %w", err)
	}

	return comments, nil
}

func (r *ConsumerRepo) UpdateStoreCommentVotes(ctx context.Context, storeID string, commentID string, amount string) (dto.StoreInfo, error) {
	// Update helpful votes
	queryUpdate := `
		UPDATE comments
		SET helpful_votes = helpful_votes + $1
		WHERE store_id = $2 AND comment_id = $3;
	`

	amt, err := strconv.Atoi(amount)
	if err != nil {
		return dto.StoreInfo{}, fmt.Errorf("invalid amount: %w", err)
	}

	_, err = r.DB.conn.Exec(ctx, queryUpdate, amt, storeID, commentID)
	if err != nil {
		return dto.StoreInfo{}, fmt.Errorf("failed to update comment votes: %w", err)
	}

	// Return updated store info
	storeInfo, err := r.GetStoreInfo(ctx, storeID)
	if err != nil {
		return dto.StoreInfo{}, fmt.Errorf("failed to get updated store info: %w", err)
	}

	return storeInfo, nil
}
