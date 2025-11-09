package db

import (
	"context"
	"strconv"

	"backend/internal/consumer-service/core/domain/dto"
	"backend/internal/mylogger"

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

func (r *ConsumerRepo) GetProductInfo(ctx context.Context, productID string) (dto.ProductInfo, error) {
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
		return dto.ProductInfo{}, err
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
		return dto.ProductInfo{}, err
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
			return dto.ProductInfo{}, err
		}
		stores = append(stores, s)
	}

	if err = rows.Err(); err != nil {
		log.Error().
			Err(err).
			Str("query", "GetProductInfo - stores").
			Msg("row iteration error")
		return dto.ProductInfo{}, err
	}

	log.Debug().
		Str("product_id", productID).
		Int("store_count", len(stores)).
		Msg("successfully fetched product info")

	return dto.ProductInfo{
		Product: product,
		Stores:  stores,
	}, nil
}
