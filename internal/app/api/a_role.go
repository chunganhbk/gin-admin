package api

import (
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/internal/app/services"
	"github.com/chunganhbk/gin-go/pkg/app"
	"github.com/gin-gonic/gin"
)

type Role struct {
	RoleService services.IRoleService
}

func NewRole(roleService services.IRoleService) *Role {
	return &Role{RoleService: roleService}
}

// Query menu
func (r *Role) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RoleQueryParam
	if err := app.ParseQuery(c, &params); err != nil {
		app.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := r.RoleService.Query(ctx, params, schema.RoleQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("order_number", schema.OrderByDESC)),
	})
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResPage(c, result.Data, result.PageResult)
}

// Query Select role
func (r *Role) QuerySelect(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RoleQueryParam
	if err := app.ParseQuery(c, &params); err != nil {
		app.ResError(c, err)
		return
	}

	result, err := r.RoleService.Query(ctx, params, schema.RoleQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("order_number", schema.OrderByDESC)),
	})
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResList(c, result.Data)
}

// Get role detail
func (r *Role) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := r.RoleService.Get(ctx, c.Param("id"))
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResSuccess(c, item)
}

// Create menu
func (a *Role) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Role
	if err := app.ParseJSON(c, &item); err != nil {
		app.ResError(c, err)
		return
	}

	item.Creator = app.GetUserID(c)
	result, err := a.RoleService.Create(ctx, item)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResSuccess(c, result)
}

// Update role
func (r *Role) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Role
	if err := app.ParseJSON(c, &item); err != nil {
		app.ResError(c, err)
		return
	}

	err := r.RoleService.Update(ctx, c.Param("id"), item)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResOK(c)
}

// Delete  role
func (r *Role) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := r.RoleService.Delete(ctx, c.Param("id"))
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResOK(c)
}

// Enable role
func (r *Role) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := r.RoleService.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResOK(c)
}

// Disable role
func (r *Role) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := r.RoleService.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResOK(c)
}
