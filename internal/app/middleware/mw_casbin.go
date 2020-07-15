package middleware

import (
	"github.com/chunganhbk/gin-go/internal/app/config"
	"github.com/casbin/casbin/v2"
	"github.com/chunganhbk/gin-go/pkg/app"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Casbin Middleware
func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) gin.HandlerFunc {
	cfg := config.C.Casbin
	if !cfg.Enable {
		return EmptyMiddleware()
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}
		p := c.Request.URL.Path
		m := c.Request.Method
		if b, err := enforcer.Enforce(app.GetUserID(c), p, m); err != nil {
			app.ResError(c, errors.WithStack(err))
			return
		} else if !b {
			app.ResError(c, app.NoPermissionResponse())
			return
		}
		c.Next()
	}
}
