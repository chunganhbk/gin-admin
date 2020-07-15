package gorm

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetMenuActionResourceDB
func GetMenuActionResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(MenuActionResource))
}

// SchemaMenuActionResource
type SchemaMenuActionResource schema.MenuActionResource

// ToMenuActionResource
func (a SchemaMenuActionResource) ToMenuActionResource() *MenuActionResource {
	item := new(MenuActionResource)
	util.StructMapToStruct(a, item)
	return item
}

// MenuActionResource
type MenuActionResource struct {
	Model
	ActionID string `gorm:"column:action_id;size:36;index;default:'';not null;"`
	Method   string `gorm:"column:method;size:100;default:'';not null;"`
	Path     string `gorm:"column:path;size:100;default:'';not null;"`
}



// To Schema MenuAction Resource
func (a MenuActionResource) ToSchemaMenuActionResource() *schema.MenuActionResource {
	item := new(schema.MenuActionResource)
	util.StructMapToStruct(a, item)
	return item
}

// Menu Action Resources
type MenuActionResources []*MenuActionResource

// To Schema MenuAction Resources
func (a MenuActionResources) ToSchemaMenuActionResources() []*schema.MenuActionResource {
	list := make([]*schema.MenuActionResource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuActionResource()
	}
	return list
}
