package gorm

import (
	"context"

	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetRoleDB
func GetRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(Role))
}

// SchemaRole
type SchemaRole schema.Role

// ToRole
func (a SchemaRole) ToRole() *Role {
	item := new(Role)
	util.StructMapToStruct(a, item)
	return item
}

// Role
type Role struct {
	Model
	Name    string  `gorm:"column:name;size:100;default:'';not null;"`
	Order   int     `gorm:"column:order_number;index;default:0;not null;"`
	Memo    *string `gorm:"column:memo;size:1024;"`
	Status  int     `gorm:"column:status;default:0;not null;"`
	Creator string  `gorm:"column:creator;size:36;"`
}

// To Schema Role
func (a Role) ToSchemaRole() *schema.Role {
	item := new(schema.Role)
	util.StructMapToStruct(a, item)
	return item
}

// Roles
type Roles []*Role

// ToSchemaRoles
func (a Roles) ToSchemaRoles() []*schema.Role {
	list := make([]*schema.Role, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRole()
	}
	return list
}
