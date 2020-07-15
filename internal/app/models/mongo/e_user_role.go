package mongo

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUserRoleCollection UserRole
func GetUserRoleCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return GetCollection(ctx, cli, UserRole{})
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
	Model  `bson:",inline"`
	UserID string `bson:"user_id"`
	RoleID string `bson:"role_id"`
}

// CollectionName
func (a UserRole) CollectionName() string {
	return a.Model.CollectionName("user_role")
}

// CreateIndexes
func (a UserRole) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"user_id": 1}},
		{Keys: bson.M{"role_id": 1}},
	})
}

// ToSchemaUserRole
func (a UserRole) ToSchemaUserRole() *schema.UserRole {
	item := new(schema.UserRole)
	util.StructMapToStruct(a, item)
	return item
}

// UserRoles
type UserRoles []*UserRole

// ToSchemaUserRoles
func (a UserRoles) ToSchemaUserRoles() []*schema.UserRole {
	list := make([]*schema.UserRole, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUserRole()
	}
	return list
}
