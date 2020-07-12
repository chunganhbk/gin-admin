package middleware

import (
	"github.com/chunganhbk/gin-go/internal/app/config"
	"github.com/chunganhbk/gin-go/internal/app/ginplus"
	"github.com/chunganhbk/gin-go/internal/app/icontext"
	"github.com/chunganhbk/gin-go/pkg/app"
	"github.com/chunganhbk/gin-go/pkg/auth"
	"github.com/chunganhbk/gin-go/pkg/errors"
	"github.com/chunganhbk/gin-go/pkg/jwt"
	"github.com/chunganhbk/gin-go/pkg/logger"
	"github.com/gin-gonic/gin"
)

func wrapUserAuthContext(c *gin.Context, userID string) {
	ginplus.SetUserID(c, userID)
	ctx := icontext.NewUserID(c.Request.Context(), userID)
	ctx = logger.NewUserIDContext(ctx, userID)
	c.Request = c.Request.WithContext(ctx)
}

// UserAuthMiddleware 用户授权中间件
func UserAuthMiddleware(a jwt.IJWTAuth, skippers ...SkipperFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userID, err := a.ParseUserID(c.Request.Context(), app.GetToken(c))
		if err != nil {
			if err == auth.ErrInvalidToken {
				if config.C.IsDebugMode() {
					wrapUserAuthContext(c, config.C.Root.UserName)
					c.Next()
					return
				}
				ginplus.ResError(c, errors.ErrInvalidToken)
				return
			}
			ginplus.ResError(c, errors.WithStack(err))
			return
		}

		wrapUserAuthContext(c, userID)
		c.Next()
	}
}
