package middleware

import (
	"github.com/chunganhbk/gin-go/internal/app/config"
	"github.com/chunganhbk/gin-go/internal/app/ginplus"
	"github.com/chunganhbk/gin-go/pkg/errors"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware casbin中间件
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
		if b, err := enforcer.Enforce(ginplus.GetUserID(c), p, m); err != nil {
			ginplus.ResError(c, errors.WithStack(err))
			return
		} else if !b {
			ginplus.ResError(c, errors.ErrNoPerm)
			return
		}
		c.Next()
	}
}
