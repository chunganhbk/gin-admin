package impl

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/repositories"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/app"
	"github.com/chunganhbk/gin-go/pkg/jwt"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/pkg/errors"
)



// Login
type AuthService struct {
	Jwt            jwt.IJWTAuth
	UserRp       repositories.IUser
}
func NewAuthService(jwt jwt.IJWTAuth, userRp repositories.IUser) *AuthService{
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
		return nil, app.NewResponse(app.ERROR_LOGIN_FAILED, app.INVALID_PARAMS, nil)
	}

	item := result.Data[0]
	if util.ComparePasswords(password, item.Password) {
		return nil, app.NewResponse(app.ERROR_LOGIN_FAILED, app.INVALID_PARAMS, nil)
	} else if item.Status != 1 {
		return nil, app.NewResponse(app.ERROR_USER_DISABLED, app.INVALID_PARAMS, nil)
	}

	return item, nil
}

// Generate Token
func (a *AuthService) GenerateToken( userID string) (*schema.LoginTokenInfo, error) {
	tokenInfo, err := a.Jwt.GenerateToken(userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	item := &schema.LoginTokenInfo{
		AccessToken: tokenInfo.GetAccessToken(),
		TokenType:   tokenInfo.GetTokenType(),
		ExpiresAt:   tokenInfo.GetExpiresAt(),
	}
	return item, nil
}
