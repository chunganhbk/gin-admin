package gorm

import (
	"context"

	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetMenuDB
func GetMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(Menu))
}

// SchemaMenu
type SchemaMenu schema.Menu

// ToMenu
func (a SchemaMenu) ToMenu() *Menu {
	item := new(Menu)
	util.StructMapToStruct(a, item)
	return item
}

// Menu
type Menu struct {
	Model
	Name       string  `gorm:"column:name;size:50;index;default:'';not null;"`
	Order      int     `gorm:"column:order_number;index;default:0;not null;"`
	Icon       *string `gorm:"column:icon;size:255;"`
	Router     *string `gorm:"column:router;size:255;"`
	ParentID   *string `gorm:"column:parent_id;size:36;index;"`
	ParentPath *string `gorm:"column:parent_path;size:518;index;"`
	ShowStatus int     `gorm:"column:show_status;index;default:0;not null;"`
	Status     int     `gorm:"column:status;index;default:0;not null;"`
	Memo       *string `gorm:"column:memo;size:1024;"`
	Creator    string  `gorm:"column:creator;size:36;"`
}

// ToSchemaMenu
func (a Menu) ToSchemaMenu() *schema.Menu {
	item := new(schema.Menu)
	util.StructMapToStruct(a, item)
	return item
}

// Menus
type Menus []*Menu

// ToSchemaMenus
func (a Menus) ToSchemaMenus() []*schema.Menu {
	list := make([]*schema.Menu, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenu()
	}
	return list
}
