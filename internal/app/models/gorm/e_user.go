package gorm

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/jinzhu/gorm"
)

// GetUserDB
func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(User))
}

// SchemaUser
type SchemaUser schema.User

// ToUser
func (a SchemaUser) ToUser() *User {
	item := new(User)
	util.StructMapToStruct(a, item)
	return item
}

// User
type User struct {
	Model
	UserName string  `gorm:"column:user_name;size:64;index;default:'';not null;"`
	FullName string  `gorm:"column:full_name;size:64;index;default:'';not null;"`
	Password string  `gorm:"column:password;size:40;default:'';not null;"`
	Email    *string `gorm:"column:email;size:255;index;"`
	Phone    *string `gorm:"column:phone;size:20;index;"`
	Status   int     `gorm:"column:status;index;default:0;not null;"`
	Creator  string  `gorm:"column:creator;size:36;"`
}

// ToSchemaUser
func (a User) ToSchemaUser() *schema.User {
	item := new(schema.User)
	util.StructMapToStruct(a, item)
	return item
}

// Users
type Users []*User

// ToSchemaUsers
func (a Users) ToSchemaUsers() []*schema.User {
	list := make([]*schema.User, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUser()
	}
	return list
}
