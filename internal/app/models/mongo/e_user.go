package mongo

import (
	"context"

	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUserCollection User
func GetUserCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return GetCollection(ctx, cli, User{})
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
	Model            `bson:",inline"`
	UserName  string `bson:"user_name"`
	FullName  string `bson:"full_name"`
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Password  string `bson:"password"`
	Email     string `bson:"email"`
	Phone     string `bson:"phone"`
	Status    int    `bson:"status"`
	Creator   string `bson:"creator"`
}

// CollectionName
func (a User) CollectionName() string {
	return a.Model.CollectionName("user")
}

// CreateIndexes
func (a User) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"email": 1}},
		{Keys: bson.M{"status": 1}},
	})
}

// To Schema User
func (a User) ToSchemaUser() *schema.User {
	item := new(schema.User)
	util.StructMapToStruct(a, item)
	return item
}

// Users
type Users []*User

// To Schema Users
func (a Users) ToSchemaUsers() []*schema.User {
	list := make([]*schema.User, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUser()
	}
	return list
}
