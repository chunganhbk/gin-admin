package services

import (
	"context"

	"github.com/chunganhbk/gin-go/internal/app/schema"
)

// IUserService
type IUserService interface {
	InitData(ctx context.Context) error

	GetLoginInfo(ctx context.Context, userID string) (*schema.UserLoginInfo, error)

	QueryUserMenuTree(ctx context.Context, userID string) (schema.MenuTrees, error)

	Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error)

	QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error)

	Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error)

	Create(ctx context.Context, item schema.User) (*schema.IDResult, error)

	Register(ctx context.Context, item schema.RegisterUser) (*schema.IDResult, error)

	Update(ctx context.Context, id string, item schema.User) error

	Delete(ctx context.Context, id string) error

	UpdateStatus(ctx context.Context, id string, status int) error

	ChangePassword(ctx context.Context, userID string, params schema.UpdatePasswordParam) error
}
