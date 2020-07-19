package router

import (
	"github.com/casbin/casbin/v2"
	"github.com/chunganhbk/gin-go/internal/app/api"
	"github.com/chunganhbk/gin-go/internal/app/middleware"
	"github.com/chunganhbk/gin-go/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

// RegisterAPI register api group router
func RegisterAPI(r *gin.Engine, container *dig.Container) error {
	err := api.Inject(container)
	if err != nil {
		return err
	}
	return container.Invoke(func(
		jwt jwt.IJWTAuth,
		e *casbin.SyncedEnforcer,
		cAuth *api.Auth,
		cMenu *api.Menu,
		cRole *api.Role,
		cUser *api.User,
	) error {
		g := r.Group("/api")

		g.Use(middleware.UserAuthMiddleware(jwt,
			middleware.AllowPathPrefixSkipper("/api/v1/auth", "/api/v1/pub/refresh-token"),
		))

		g.Use(middleware.CasbinMiddleware(e,
			middleware.AllowPathPrefixSkipper("/api/v1/pub", "/api/v1/auth"),
		))

		v1 := g.Group("/v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/login", cAuth.Login)
				auth.POST("/register", cAuth.Register)

			}
			gCurrent := v1.Group("/pub") //public api
			{
				gCurrent.PUT("/password", cUser.ChangePassword)
				gCurrent.GET("/user", cUser.GetUserInfo)
				gCurrent.GET("/menu-tree", cUser.QueryUserMenuTree)
				gCurrent.POST("/refresh-token", cAuth.RefreshToken)
			}

			gMenu := v1.Group("menus")
			{
				gMenu.GET("", cMenu.Query)
				gMenu.GET(":id", cMenu.Get)
				gMenu.POST("", cMenu.Create)
				gMenu.PUT(":id", cMenu.Update)
				gMenu.DELETE(":id", cMenu.Delete)
				gMenu.PATCH(":id/enable", cMenu.Enable)
				gMenu.PATCH(":id/disable", cMenu.Disable)
			}
			v1.GET("/menus.tree", cMenu.QueryTree)

			gRole := v1.Group("roles")
			{
				gRole.GET("", cRole.Query)
				gRole.GET(":id", cRole.Get)
				gRole.POST("", cRole.Create)
				gRole.PUT(":id", cRole.Update)
				gRole.DELETE(":id", cRole.Delete)
				gRole.PATCH(":id/enable", cRole.Enable)
				gRole.PATCH(":id/disable", cRole.Disable)
			}
			v1.GET("/roles.select", cRole.QuerySelect)

			gUser := v1.Group("users")
			{
				gUser.GET("", cUser.Query)
				gUser.GET(":id", cUser.Get)
				gUser.POST("", cUser.Create)
				gUser.PUT(":id", cUser.Update)
				gUser.DELETE(":id", cUser.Delete)
				gUser.PATCH(":id/enable", cUser.Enable)
				gUser.PATCH(":id/disable", cUser.Disable)
			}
		}
		return nil
	})
}
