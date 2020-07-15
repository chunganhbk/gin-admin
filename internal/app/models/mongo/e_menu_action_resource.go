package mongo

import (
	"context"

	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetMenuActionResourceCollection
func GetMenuActionResourceCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return GetCollection(ctx, cli, MenuActionResource{})
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
	Model    `bson:",inline"`
	ActionID string `bson:"action_id"`
	Method   string `bson:"method"`
	Path     string `bson:"path"`
}

// CollectionName
func (a MenuActionResource) CollectionName() string {
	return a.Model.CollectionName("menu_action_resource")
}

// CreateIndexes
func (a MenuActionResource) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"action_id": 1}},
	})
}

// ToSchemaMenuActionResource
func (a MenuActionResource) ToSchemaMenuActionResource() *schema.MenuActionResource {
	item := new(schema.MenuActionResource)
	util.StructMapToStruct(a, item)
	return item
}

// MenuActionResources
type MenuActionResources []*MenuActionResource

// ToSchema MenuActionResources
func (a MenuActionResources) ToSchemaMenuActionResources() []*schema.MenuActionResource {
	list := make([]*schema.MenuActionResource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuActionResource()
	}
	return list
}
