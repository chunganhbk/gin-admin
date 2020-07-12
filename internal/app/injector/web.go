package injector

import (
	"github.com/chunganhbk/gin-go/internal/app/config"
	"github.com/chunganhbk/gin-go/internal/app/middleware"
	"github.com/chunganhbk/gin-go/internal/app/router"
	"github.com/LyricTian/gzip"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

// InitGinEngine 初始化gin引擎
func InitGinEngine(r router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())
	prefixes := r.Prefixes()
	// Access logger
	app.Use(middleware.LoggerMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))


	// CORS
	if config.C.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	// GZIP
	if config.C.GZIP.Enable {
		app.Use(gzip.Gzip(gzip.BestCompression,
			gzip.WithExcludedExtensions(config.C.GZIP.ExcludedExtentions),
			gzip.WithExcludedPaths(config.C.GZIP.ExcludedPaths),
		))
	}

	// Router register
	r.Register(app)

	// Swagger
	if config.C.Swagger {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Website
	if dir := config.C.WWW; dir != "" {
		app.Use(middleware.WWWMiddleware(dir, middleware.AllowPathPrefixSkipper(prefixes...)))
	}

	return app
}
