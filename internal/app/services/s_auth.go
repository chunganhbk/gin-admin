package services

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/schema"
)

// ILogin
type IAuthService interface {
	Verify(ctx context.Context, Email string, password string) (*schema.User, error)

	GenerateToken(userID string) (*schema.LoginTokenInfo, error)
}
