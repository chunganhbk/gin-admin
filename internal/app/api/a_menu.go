package api

import (
	"github.com/chunganhbk/gin-go/internal/app/services"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/app"
	"github.com/gin-gonic/gin"

)

// Menu
type Menu struct {
	MenuService services.IMenuService
}

func NewMenu( menuService services.IMenuService) *Menu{
	return &Menu{menuService}
}
// Query menu
func (m *Menu) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := app.ParseQuery(c, &params); err != nil {
		app.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := m.MenuService.Query(ctx, params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResPage(c, result.Data, result.PageResult)
}

// Query Tree  menu
func (m *Menu) QueryTree(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := app.ParseQuery(c, &params); err != nil {
		app.ResError(c, err)
		return
	}

	result, err := m.MenuService.Query(ctx, params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("order", schema.OrderByDESC)),
	})
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResList(c, result.Data.ToTree())
}

// Get menu
func (a *Menu) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.MenuService.Get(ctx, c.Param("id"))
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResSuccess(c, item)
}

// Create menu
func (m *Menu) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Menu
	if err := app.ParseJSON(c, &item); err != nil {
		app.ResError(c, err)
		return
	}

	item.Creator = app.GetUserID(c)
	result, err := m.MenuService.Create(ctx, item)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResSuccess(c, result)
}

// Update menu
func (m *Menu) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Menu
	if err := app.ParseJSON(c, &item); err != nil {
		app.ResError(c, err)
		return
	}

	err := m.MenuService.Update(ctx, c.Param("id"), item)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResSuccess(c, nil)
}

// Delete menu
func (m *Menu) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := m.MenuService.Delete(ctx, c.Param("id"))
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResOK(c)
}

// Enable menu
func (m *Menu) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := m.MenuService.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResOK(c)
}

// Disable menu
func (m *Menu) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := m.MenuService.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		app.ResError(c, err)
		return
	}
	app.ResOK(c)
}
