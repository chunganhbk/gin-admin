package api

import (
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/internal/app/services"
	"github.com/chunganhbk/gin-go/pkg/app"
	"github.com/chunganhbk/gin-go/pkg/errors"
	"github.com/gin-gonic/gin"
	"strings"
)

// User
type User struct {
	UserService services.IUserService
}

func NewUser(userService services.IUserService) *User {
	return &User{userService}
}

// Get User Info
func (u *User) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := u.UserService.GetLoginInfo(ctx, app.GetUserID(c))
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResSuccess(c, info)
}

// Query User MenuTree
func (u *User) QueryUserMenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := u.UserService.QueryUserMenuTree(ctx, app.GetUserID(c))
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResList(c, menus)
}

// Update Password  current user
func (u *User) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.UpdatePasswordParam
	if err := app.ParseJSON(c, &item); err != nil {
		app.ResError(c, err)
		return
	}

	err := u.UserService.ChangePassword(ctx, app.GetUserID(c), item)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResOK(c)
}

// Query user list
func (a *User) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.UserQueryParam
	if err := app.ParseQuery(c, &params); err != nil {
		app.ResError(c, err)
		return
	}
	if v := c.Query("roleIDs"); v != "" {
		params.RoleIDs = strings.Split(v, ",")
	}

	params.Pagination = true
	result, err := a.UserService.QueryShow(ctx, params)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResPage(c, result.Data, result.PageResult)
}

// Get user info
func (u *User) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := u.UserService.Get(ctx, c.Param("id"))
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResSuccess(c, item.CleanSecure())
}

// Create user
func (a *User) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := app.ParseJSON(c, &item); err != nil {
		app.ResError(c, err)
		return
	} else if item.Password == "" {
		app.ResError(c, errors.New400Response(errors.ERROR_NOT_EXIST_USER))
		return
	}

	item.Creator = app.GetUserID(c)
	result, err := a.UserService.Create(ctx, item)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResSuccess(c, result)
}

// Update user
func (u *User) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := app.ParseJSON(c, &item); err != nil {
		app.ResError(c, err)
		return
	}

	err := u.UserService.Update(ctx, c.Param("id"), item)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResOK(c)
}

// Delete user
func (u *User) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := u.UserService.Delete(ctx, c.Param("id"))
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResOK(c)
}

// Enable  status user
func (u *User) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := u.UserService.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResOK(c)
}

// Disable status user
func (u *User) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := u.UserService.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResOK(c)
}
