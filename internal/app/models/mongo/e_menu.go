package mongo

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetMenuCollection
func GetMenuCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return GetCollection(ctx, cli, Menu{})
}

// SchemaMenu 菜单对象
type SchemaMenu schema.Menu

// ToMenu 转换为菜单实体
func (a SchemaMenu) ToMenu() *Menu {
	item := new(Menu)
	util.StructMapToStruct(a, item)
	return item
}

// Menu 菜单实体
type Menu struct {
	Model      `bson:",inline"`
	Name       string `bson:"name"`
	Order      int    `bson:"order"`
	Icon       string `bson:"icon"`
	Router     string `bson:"router"`
	ParentID   string `bson:"parent_id"`
	ParentPath string `bson:"parent_path"`
	ShowStatus int    `bson:"show_status"`
	Status     int    `bson:"status"`
	Memo       string `bson:"memo"`
	Creator    string `bson:"creator"`
}

// CollectionName
func (a Menu) CollectionName() string {
	return a.Model.CollectionName("menu")
}

// CreateIndexes
func (a Menu) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"sequence": -1}},
		{Keys: bson.M{"parent_id": 1}},
		{Keys: bson.M{"parent_path": 1}},
		{Keys: bson.M{"show_status": 1}},
		{Keys: bson.M{"status": 1}},
	})
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
