package db

import (
	"context"

	"backend/internal/business-service/core/domain/dto"
	"backend/internal/mylogger"
)

type BusinessRepo struct {
	ctx   context.Context
	mylog mylogger.Logger
	DB    *DB
}

func NewBusinessRepo(ctx context.Context, db *DB, mylog mylogger.Logger) *BusinessRepo {
	return &BusinessRepo{
		ctx:   ctx,
		mylog: mylog,
		DB:    db,
	}
}

func (r *BusinessRepo) GetStores(ctx context.Context, userId string) ([]dto.Store, error) {
	query := `
		SELECT 
			s.store_id,
			s.name,
			s.address,
			s.photo,
			s.description,
			COALESCE(c.latitude, 0) AS latitude,
			COALESCE(c.longitude, 0) AS longitude
		FROM stores s
		LEFT JOIN coords c ON s.coord = c.coord_id
		WHERE s.owner_id = $1;
	`

	rows, err := r.DB.conn.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []dto.Store

	for rows.Next() {
		var s dto.Store
		err := rows.Scan(
			&s.StoreID,
			&s.Name,
			&s.Address,
			&s.Photo,
			&s.Description,
			&s.Latitude,
			&s.Longitude,
		)
		if err != nil {
			return nil, err
		}
		stores = append(stores, s)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return stores, nil
}

func (r *BusinessRepo) AddStore(ctx context.Context, userId string, store dto.Store) error {
	tx, err := r.DB.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1️⃣ Insert the store first, returning the generated store_id
	var storeID string
	queryStore := `
        INSERT INTO stores (owner_id, name, address, photo, description)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING store_id;
    `
	err = tx.QueryRow(ctx, queryStore,
		userId,
		store.Name,
		store.Address,
		store.Photo,
		store.Description,
	).Scan(&storeID)
	if err != nil {
		return err
	}

	var coordID *string

	// 2️⃣ If coordinates exist, insert into coords table using the new store_id
	if store.Latitude != 0 && store.Longitude != 0 {
		queryCoord := `
            INSERT INTO coords (entity_id, entity_type, latitude, longitude, address)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING coord_id;
        `
		var newCoordID string
		err = tx.QueryRow(ctx, queryCoord,
			storeID, "store", store.Latitude, store.Longitude, store.Address,
		).Scan(&newCoordID)
		if err != nil {
			return err
		}
		coordID = &newCoordID
	}

	// 3️⃣ Update store to link coord_id (if it was created)
	if coordID != nil {
		queryUpdate := `
            UPDATE stores
            SET coord = $1
            WHERE store_id = $2;
        `
		_, err = tx.Exec(ctx, queryUpdate, coordID, storeID)
		if err != nil {
			return err
		}
	}

	// ✅ Commit transaction
	return tx.Commit(ctx)
}

func (r *BusinessRepo) UpdateStore(ctx context.Context, store dto.Store) error {
	query := `
		UPDATE stores
		SET name = $1,
			address = $2,
			photo = $3,
			description = $4,
			updated_at = NOW()
		WHERE store_id = $5;
	`

	_, err := r.DB.conn.Exec(ctx, query,
		store.Name,
		store.Address,
		store.Photo,
		store.Description,
		store.StoreID,
	)
	return err
}

func (r *BusinessRepo) DeleteStore(ctx context.Context, storeId string) error {
	tx, err := r.DB.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Remove products linked to this store
	_, err = tx.Exec(ctx, `
        DELETE FROM products WHERE store_id = $1;
    `, storeId)
	if err != nil {
		return err
	}

	// 2. Remove the store itself
	_, err = tx.Exec(ctx, `
        DELETE FROM stores WHERE store_id = $1;
    `, storeId)
	if err != nil {
		return err
	}

	// 3. Remove coords linked to this store
	_, err = tx.Exec(ctx, `
        DELETE FROM coords WHERE entity_id = $1 AND entity_type = 'store';
    `, storeId)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *BusinessRepo) GetProducts(ctx context.Context, storeId string) ([]dto.Product, error) {
	query := `
		SELECT product_id, name, description, photo, price
		FROM products
		WHERE store_id = $1;
	`

	rows, err := r.DB.conn.Query(ctx, query, storeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []dto.Product
	for rows.Next() {
		var p dto.Product
		err := rows.Scan(
			&p.ProductID,
			&p.Name,
			&p.Description,
			&p.Photo,
			&p.Price,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return products, nil
}

func (r *BusinessRepo) AddProduct(ctx context.Context, storeId string, product dto.Product) error {
	query := `
		INSERT INTO products (store_id, name, description, photo, price, in_stock)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := r.DB.conn.Exec(ctx, query,
		storeId, // replace when you have storeID available
		product.Name,
		product.Description,
		product.Photo,
		product.Price,
		product.InStock,
	)
	return err
}

func (r *BusinessRepo) UpdateProduct(ctx context.Context, product dto.Product) error {
	query := `
		UPDATE products
		SET name = $1,
			description = $2,
			photo = $3,
			price = $4,
			updated_at = NOW(),
			in_stock = $5
		WHERE product_id = $6;
	`

	_, err := r.DB.conn.Exec(ctx, query,
		product.Name,
		product.Description,
		product.Photo,
		product.Price,
		product.InStock,
		product.ProductID,
	)
	return err
}

func (r *BusinessRepo) DeleteProduct(ctx context.Context, productId string) error {
	_, err := r.DB.conn.Exec(ctx, `DELETE FROM products WHERE product_id = $1;`, productId)
	return err
}
