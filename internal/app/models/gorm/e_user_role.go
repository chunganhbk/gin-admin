package gorm

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetUserRoleDB
func GetUserRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(UserRole))
}

// SchemaUserRole
type SchemaUserRole schema.UserRole

// ToUserRole
func (a SchemaUserRole) ToUserRole() *UserRole {
	item := new(UserRole)
	util.StructMapToStruct(a, item)
	return item
}

// UserRole
type UserRole struct {
	Model
	UserID string `gorm:"column:user_id;size:36;index;default:'';not null;"`
	RoleID string `gorm:"column:role_id;size:36;index;default:'';not null;"`
}

// To Schema UserRole
func (a UserRole) ToSchemaUserRole() *schema.UserRole {
	item := new(schema.UserRole)
	util.StructMapToStruct(a, item)
	return item
}

// UserRoles
type UserRoles []*UserRole

// To Schema UserRoles
func (a UserRoles) ToSchemaUserRoles() []*schema.UserRole {
	list := make([]*schema.UserRole, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUserRole()
	}
	return list
}
