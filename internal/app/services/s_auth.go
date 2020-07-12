package services

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/schema"
)

// ILogin
type IAuthService interface {


	Verify(ctx context.Context, userName, password string) (*schema.User, error)

	GenerateToken(ctx context.Context, userID string) (*schema.LoginTokenInfo, error)

}
