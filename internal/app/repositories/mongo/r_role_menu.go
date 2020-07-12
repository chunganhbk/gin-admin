package mongo

import (
	"context"
	"github.com/pkg/errors"
	"time"
	entity "github.com/chunganhbk/gin-go/internal/app/models/mongo"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Role Menu
type RoleMenu struct {
	Client *mongo.Client
}
func NewRoleMenu(client *mongo.Client) *RoleMenu{
	return &RoleMenu{client}
}
func (a *RoleMenu) getQueryOption(opts ...schema.RoleMenuQueryOptions) schema.RoleMenuQueryOptions {
	var opt schema.RoleMenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}


func (a *RoleMenu) Query(ctx context.Context, params schema.RoleMenuQueryParam, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenuQueryResult, error) {
	opt := a.getQueryOption(opts...)

	c := entity.GetRoleMenuCollection(ctx, a.Client)
	filter := DefaultFilter(ctx)
	roleIDs := params.RoleIDs
	if v := params.RoleID; v != "" {
		roleIDs = append(roleIDs, v)
	}
	if v := roleIDs; len(v) > 0 {
		filter = append(filter, Filter("role_id", bson.M{"$in": v}))
	}
	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByDESC))

	var list entity.RoleMenus
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleMenuQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoleMenus(),
	}

	return qr, nil
}

func (a *RoleMenu) Get(ctx context.Context, id string, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenu, error) {
	c := entity.GetRoleMenuCollection(ctx, a.Client)
	filter := DefaultFilter(ctx, Filter("_id", id))
	var item entity.RoleMenu
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaRoleMenu(), nil
}

// Create
func (a *RoleMenu) Create(ctx context.Context, item schema.RoleMenu) error {
	eitem := entity.SchemaRoleMenu(item).ToRoleMenu()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetRoleMenuCollection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	return errors.WithStack(err)
}

func (a *RoleMenu) Update(ctx context.Context, id string, item schema.RoleMenu) error {
	eitem := entity.SchemaRoleMenu(item).ToRoleMenu()
	eitem.UpdatedAt = time.Now()
	c := entity.GetRoleMenuCollection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", id)), eitem)
	return errors.WithStack(err)
}


func (a *RoleMenu) Delete(ctx context.Context, id string) error {
	c := entity.GetRoleMenuCollection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", id)))
	return errors.WithStack(err)
}

// DeleteByRoleID
func (a *RoleMenu) DeleteByRoleID(ctx context.Context, roleID string) error {
	c := entity.GetRoleMenuCollection(ctx, a.Client)
	err := DeleteMany(ctx, c, DefaultFilter(ctx, Filter("role_id", roleID)))
	return errors.WithStack(err)
}
