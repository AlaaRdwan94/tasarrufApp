package payment

import (
	"context"
)

// Payment represents a transaction.
type Payment interface {
	Submit(ctx context.Context) (string, error)
	Cancel(ctx context.Context) (string, error)
}
