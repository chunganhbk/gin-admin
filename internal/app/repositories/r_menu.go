package repositories

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/schema"
)

// IMenu
type IMenu interface {

	Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error)

	Get(ctx context.Context, id string, opts ...schema.MenuQueryOptions) (*schema.Menu, error)

	Create(ctx context.Context, item schema.Menu) error

	Update(ctx context.Context, id string, item schema.Menu) error

	Delete(ctx context.Context, id string) error

	UpdateParentPath(ctx context.Context, id, parentPath string) error

	UpdateStatus(ctx context.Context, id string, status int) error
}
