package middleware

import (
	"time"
	"github.com/chunganhbk/gin-go/internal/app/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS Middleware
func CORSMiddleware() gin.HandlerFunc {
	cfg := config.C.CORS
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Second * time.Duration(cfg.MaxAge),
	})
}
