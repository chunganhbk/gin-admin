package gorm

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetRoleMenuDB
func GetRoleMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(RoleMenu))
}

// SchemaRoleMenu
type SchemaRoleMenu schema.RoleMenu

// ToRoleMenu
func (a SchemaRoleMenu) ToRoleMenu() *RoleMenu {
	item := new(RoleMenu)
	util.StructMapToStruct(a, item)
	return item
}

// RoleMenu
type RoleMenu struct {
	ID       string `gorm:"column:id;primary_key;size:36;"`
	RoleID   string `gorm:"column:role_id;size:36;index;default:'';not null;"`
	MenuID   string `gorm:"column:menu_id;size:36;index;default:'';not null;"`
	ActionID string `gorm:"column:action_id;size:36;default:'';not null;"`
}

// To Schema RoleMenu
func (a RoleMenu) ToSchemaRoleMenu() *schema.RoleMenu {
	item := new(schema.RoleMenu)
	util.StructMapToStruct(a, item)
	return item
}

// Role Menus
type RoleMenus []*RoleMenu

// To Schema RoleMenus
func (a RoleMenus) ToSchemaRoleMenus() []*schema.RoleMenu {
	list := make([]*schema.RoleMenu, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRoleMenu()
	}
	return list
}
