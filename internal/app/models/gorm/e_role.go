package gorm

import (
	"context"

	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetRoleDB 获取角色存储
func GetRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(Role))
}

// SchemaRole 角色对象
type SchemaRole schema.Role

// ToRole 转换为角色实体
func (a SchemaRole) ToRole() *Role {
	item := new(Role)
	util.StructMapToStruct(a, item)
	return item
}

// Role 角色实体
type Role struct {
	Model
	Name     string  `gorm:"column:name;size:100;default:'';not null;"`
	Order int     `gorm:"column:order;index;default:0;not null;"`
	Memo     *string `gorm:"column:memo;size:1024;"`
	Status   int     `gorm:"column:status;default:0;not null;"`
	Creator  string  `gorm:"column:creator;size:36;"`
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
