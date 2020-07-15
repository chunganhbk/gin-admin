package middleware

import (
	"github.com/chunganhbk/gin-go/internal/app/icontext"
	"github.com/chunganhbk/gin-go/pkg/app"
	"github.com/chunganhbk/gin-go/pkg/jwt"
	"github.com/chunganhbk/gin-go/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func wrapUserAuthContext(c *gin.Context, userID string) {
	app.SetUserID(c, userID)
	ctx := icontext.NewUserID(c.Request.Context(), userID)
	ctx = logger.NewUserIDContext(ctx, userID)
	c.Request = c.Request.WithContext(ctx)
}

// UserAuthMiddleware
func UserAuthMiddleware(a jwt.IJWTAuth, skippers ...SkipperFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userID, err := a.ParseUserID(app.GetToken(c))
		if err != nil {
			if err == jwt.ErrInvalidToken {

				app.ResError(c, app.New400Response(app.ERROR_AUTH_CHECK_TOKEN_FAIL))
				return
			}
			app.ResError(c, errors.WithStack(err))
			return
		}
		wrapUserAuthContext(c, userID)
		c.Next()
	}
}
