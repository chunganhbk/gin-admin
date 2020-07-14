package services

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/schema"
)

// IMenuService
type IMenuService interface {

	InitData(ctx context.Context, dataFile string) error

	Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error)

	Get(ctx context.Context, id string, opts ...schema.MenuQueryOptions) (*schema.Menu, error)

	Create(ctx context.Context, item schema.Menu) (*schema.IDResult, error)

	Update(ctx context.Context, id string, item schema.Menu) error

	Delete(ctx context.Context, id string) error

	UpdateStatus(ctx context.Context, id string, status int) error
}
