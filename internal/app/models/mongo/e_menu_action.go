package mongo

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetMenuActionCollection MenuAction
func GetMenuActionCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return GetCollection(ctx, cli, MenuAction{})
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
	Model  `bson:",inline"`
	MenuID string `bson:"menu_id"` // ID
	Code   string `bson:"code"`
	Name   string `bson:"name"`
}

// CollectionName
func (a MenuAction) CollectionName() string {
	return a.Model.CollectionName("menu_action")
}

// CreateIndexes
func (a MenuAction) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"menu_id": 1}},
	})
}

// ToSchemaMenuAction
func (a MenuAction) ToSchemaMenuAction() *schema.MenuAction {
	item := new(schema.MenuAction)
	util.StructMapToStruct(a, item)
	return item
}

// MenuActions
type MenuActions []*MenuAction

// ToSchemaMenuActions
func (a MenuActions) ToSchemaMenuActions() []*schema.MenuAction {
	list := make([]*schema.MenuAction, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuAction()
	}
	return list
}
