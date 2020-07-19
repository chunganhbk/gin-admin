package mongo

import (
	"context"

	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetRoleCollection Role
func GetRoleCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return GetCollection(ctx, cli, Role{})
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
	Model   `bson:",inline"`
	Name    string `bson:"name"`
	Order   int    `bson:"order_number"`
	Memo    string `bson:"memo"`
	Status  int    `bson:"status"`
	Creator string `bson:"creator"`
}

// CollectionName role
func (a Role) CollectionName() string {
	return a.Model.CollectionName("role")
}

// CreateIndexes
func (a Role) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"order_number": -1}},
		{Keys: bson.M{"status": 1}},
	})
}

// ToSchema Role
func (a Role) ToSchemaRole() *schema.Role {
	item := new(schema.Role)
	util.StructMapToStruct(a, item)
	return item
}

// Roles
type Roles []*Role

// ToSchema Roles
func (a Roles) ToSchemaRoles() []*schema.Role {
	list := make([]*schema.Role, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRole()
	}
	return list
}
