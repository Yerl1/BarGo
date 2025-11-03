package db

import (
	"context"

	"backend/internal/mylogger"
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
