package impl

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/repositories"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/errors"
	"github.com/chunganhbk/gin-go/pkg/jwt"
	"github.com/chunganhbk/gin-go/pkg/util"
)

// Login
type AuthService struct {
	Jwt    jwt.IJWTAuth
	UserRp repositories.IUser
}

func NewAuthService(jwt jwt.IJWTAuth, userRp repositories.IUser) *AuthService {
	return &AuthService{jwt, userRp}
}

// Verify
func (a *AuthService) Verify(ctx context.Context, email string, password string) (*schema.User, error) {

	result, err := a.UserRp.Query(ctx, schema.UserQueryParam{
		Email: email,
	})

	if err != nil {
		return nil, err
	} else if len(result.Data) == 0 {
		return nil, errors.New400Response(errors.ERROR_LOGIN_FAILED)
	}

	item := result.Data[0]
	if !util.ComparePasswords(password, item.Password) {
		return nil, errors.New400Response(errors.ERROR_LOGIN_FAILED)
	} else if item.Status != 1 {
		return nil, errors.New400Response(errors.ERROR_USER_DISABLED)
	}

	return item, nil
}

// Generate Token
func (a *AuthService) GenerateToken(userID string) (*schema.LoginTokenInfo, error) {
	tokenInfo, err := a.Jwt.GenerateToken(userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	item := &schema.LoginTokenInfo{
		AccessToken:  tokenInfo.GetAccessToken(),
		RefreshToken: tokenInfo.GetRefreshToken(),
		TokenType:    tokenInfo.GetTokenType(),
		ExpiresAt:    tokenInfo.GetExpiresAt(),
	}
	return item, nil
}

//refresh token

func (a *AuthService) RefreshToken(refreshToken string) (*schema.LoginTokenInfo, error) {
	tokenInfo, err := a.Jwt.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	item := &schema.LoginTokenInfo{
		AccessToken:  tokenInfo.GetAccessToken(),
		RefreshToken: tokenInfo.GetRefreshToken(),
		TokenType:    tokenInfo.GetTokenType(),
		ExpiresAt:    tokenInfo.GetExpiresAt(),
	}
	return item, nil
}
