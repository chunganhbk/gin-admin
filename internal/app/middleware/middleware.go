package middleware

import (
	"fmt"
	"github.com/chunganhbk/gin-go/pkg/app"
	"strings"
	"github.com/gin-gonic/gin"
)

// No Method Handler
func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		app.ResError(c, app.MethodNotAllowResponse())
	}
}

// No Route Handler
func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		app.ResError(c, app.ResponseNotFound())
	}
}

// SkipperFunc
type SkipperFunc func(*gin.Context) bool

// Allow PathPrefix Skipper
func AllowPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

// Allow PathPrefix NoSkipper
func AllowPathPrefixNoSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}
		return true
	}
}

// Allow Method And PathPrefix Skipper
func AllowMethodAndPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := JoinRouter(c.Request.Method, c.Request.URL.Path)
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

// Join Router
func JoinRouter(method, path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(method), path)
}

// Skip Handler
func SkipHandler(c *gin.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}

// Empty Middleware
func EmptyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
