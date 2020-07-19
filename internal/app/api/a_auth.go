package api

import (
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/internal/app/services"
	"github.com/chunganhbk/gin-go/pkg/app"
	"github.com/chunganhbk/gin-go/pkg/errors"

	"github.com/chunganhbk/gin-go/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Login
type Auth struct {
	AuthService services.IAuthService
	UserService services.IUserService
}

func NewAuth(authService services.IAuthService, userService services.IUserService) *Auth {
	return &Auth{AuthService: authService, UserService: userService}
}

// Login
func (a *Auth) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := app.ParseJSON(c, &item); err != nil {
		app.ResError(c, err)
		return
	}

	user, err := a.AuthService.Verify(ctx, item.Email, item.Password)
	if err != nil {
		app.ResError(c, err)
		return
	}

	userID := user.ID

	app.SetUserID(c, userID)

	ctx = logger.NewUserIDContext(ctx, userID)
	tokenInfo, err := a.AuthService.GenerateToken(userID)
	if err != nil {
		app.ResError(c, err)
		return
	}

	logger.StartSpan(ctx, logger.SetSpanTitle("User login"), logger.SetSpanFuncName("Login")).Infof("Login system")
	app.ResSuccess(c, tokenInfo)
}

// Refresh Token
func (a *Auth) RefreshToken(c *gin.Context) {
	var item schema.RefreshTokenParam
	if err := app.ParseJSON(c, &item); err != nil {
		app.ResError(c, err)
		return
	}
	tokenInfo, err := a.AuthService.RefreshToken(item.RefreshToken)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResSuccess(c, tokenInfo)
}
func (a *Auth) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.RegisterUser
	if err := app.ParseJSON(c, &item); err != nil {
		app.ResError(c, err)
		return
	} else if item.Password == "" {
		app.ResError(c, errors.New400Response(errors.ERROR_PASSWORD_REQUIRED))
		return
	}
	result, err := a.UserService.Register(ctx, item)
	println("err", result, err)
	if err != nil {
		app.ResError(c, err)
		return
	}
	tokenInfo, err := a.AuthService.GenerateToken(result.ID)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResSuccess(c, tokenInfo)
}
