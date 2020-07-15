package mongo

import (
	"context"

	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetRoleMenuCollection RoleMenu
func GetRoleMenuCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return GetCollection(ctx, cli, RoleMenu{})
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
	Model    `bson:",inline"`
	RoleID   string `bson:"role_id"`
	MenuID   string `bson:"menu_id"`
	ActionID string `bson:"action_id"`
}

// CollectionName
func (a RoleMenu) CollectionName() string {
	return a.Model.CollectionName("role_menu")
}

// CreateIndexes
func (a RoleMenu) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"role_id": 1}},
		{Keys: bson.M{"menu_id": 1}},
		{Keys: bson.M{"action_id": 1}},
	})
}

// ToSchema RoleMenu
func (a RoleMenu) ToSchemaRoleMenu() *schema.RoleMenu {
	item := new(schema.RoleMenu)
	util.StructMapToStruct(a, item)
	return item
}

// RoleMenus
type RoleMenus []*RoleMenu

// ToSchema RoleMenus
func (a RoleMenus) ToSchemaRoleMenus() []*schema.RoleMenu {
	list := make([]*schema.RoleMenu, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRoleMenu()
	}
	return list
}
