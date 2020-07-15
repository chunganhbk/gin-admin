package gorm

import (
	"context"

	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetMenuActionDB
func GetMenuActionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(MenuAction))
}

// SchemaMenuAction
type SchemaMenuAction schema.MenuAction

// ToMenuAction
func (a SchemaMenuAction) ToMenuAction() *MenuAction {
	item := new(MenuAction)
	util.StructMapToStruct(a, item)
	return item
}

// MenuAction
type MenuAction struct {
	Model
	MenuID string `gorm:"column:menu_id;size:36;index;default:'';not null;"`
	Code   string `gorm:"column:code;size:100;default:'';not null;"`
	Name   string `gorm:"column:name;size:100;default:'';not null;"`
}



// ToSchema MenuAction
func (a MenuAction) ToSchemaMenuAction() *schema.MenuAction {
	item := new(schema.MenuAction)
	util.StructMapToStruct(a, item)
	return item
}

// Menu Actions
type MenuActions []*MenuAction

// ToSchema Menu Actions
func (a MenuActions) ToSchemaMenuActions() []*schema.MenuAction {
	list := make([]*schema.MenuAction, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuAction()
	}
	return list
}
