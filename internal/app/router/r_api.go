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
			middleware.AllowPathPrefixSkipper("/api/v1/auth"),
		))

		g.Use(middleware.CasbinMiddleware(e,
			middleware.AllowPathPrefixSkipper("/api/v1/auth"),
		))

		v1 := g.Group("/v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/login",cAuth.Login)
				auth.POST("/register", cAuth.Register)

			}
			gCurrent := v1.Group("current")
			{
				gCurrent.PUT("password", cUser.UpdatePassword)
				gCurrent.GET("user", cUser.GetUserInfo)
				gCurrent.GET("menutree", cUser.QueryUserMenuTree)
				gCurrent.POST("/refresh-token", cUser.RefreshToken)
			}

			gMenu := v1.Group("menus")
			{
				gMenu.GET("", a.MenuAPI.Query)
				gMenu.GET(":id", a.MenuAPI.Get)
				gMenu.POST("", a.MenuAPI.Create)
				gMenu.PUT(":id", a.MenuAPI.Update)
				gMenu.DELETE(":id", a.MenuAPI.Delete)
				gMenu.PATCH(":id/enable", a.MenuAPI.Enable)
				gMenu.PATCH(":id/disable", a.MenuAPI.Disable)
			}
			v1.GET("/menus.tree", a.MenuAPI.QueryTree)

			gRole := v1.Group("roles")
			{
				gRole.GET("", a.RoleAPI.Query)
				gRole.GET(":id", a.RoleAPI.Get)
				gRole.POST("", a.RoleAPI.Create)
				gRole.PUT(":id", a.RoleAPI.Update)
				gRole.DELETE(":id", a.RoleAPI.Delete)
				gRole.PATCH(":id/enable", a.RoleAPI.Enable)
				gRole.PATCH(":id/disable", a.RoleAPI.Disable)
			}
			v1.GET("/roles.select", a.RoleAPI.QuerySelect)

			gUser := v1.Group("users")
			{
				gUser.GET("", a.UserAPI.Query)
				gUser.GET(":id", a.UserAPI.Get)
				gUser.POST("", a.UserAPI.Create)
				gUser.PUT(":id", a.UserAPI.Update)
				gUser.DELETE(":id", a.UserAPI.Delete)
				gUser.PATCH(":id/enable", a.UserAPI.Enable)
				gUser.PATCH(":id/disable", a.UserAPI.Disable)
			}
		}
		return nil
	})
}
