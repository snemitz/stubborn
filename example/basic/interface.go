package main

import (
	"context"

	"github.com/google/uuid"
)

type BasicInterface interface {
	// Get a []byte slice or an
	Get(ctx context.Context, id uuid.UUID) ([]byte, error)
	Set(ctx context.Context, data []byte) error
	Batch(ctx context.Context, ids []uuid.UUID) ([]byte, error)
}
