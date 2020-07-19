package mongo

import (
	"context"
	entity "github.com/chunganhbk/gin-go/internal/app/models/mongo"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Role
type Role struct {
	Client *mongo.Client
}

func NewRole(client *mongo.Client) *Role {
	return &Role{client}
}
func (a *Role) getQueryOption(opts ...schema.RoleQueryOptions) schema.RoleQueryOptions {
	var opt schema.RoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query role
func (r *Role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	opt := r.getQueryOption(opts...)

	c := entity.GetRoleCollection(ctx, r.Client)
	filter := DefaultFilter(ctx)

	if v := params.IDs; len(v) > 0 {
		filter = append(filter, Filter("_id", bson.M{"$in": v}))
	}
	if v := params.Name; v != "" {
		filter = append(filter, Filter("name", v))
	}
	if v := params.UserID; v != "" {
		result, err := entity.GetUserRoleCollection(ctx, r.Client).Distinct(ctx, "role_id", DefaultFilter(ctx, Filter("user_id", v)))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		filter = append(filter, Filter("_id", bson.M{"$in": result}))
	}
	if v := params.QueryValue; v != "" {
		filter = append(filter, Filter("$or", bson.A{
			OrRegexFilter("name", v),
			OrRegexFilter("memo", v),
		}))
	}
	if v := params.Status; v > 0 {
		filter = append(filter, Filter("status", v))
	}
	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByDESC))

	var list entity.Roles
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoles(),
	}

	return qr, nil
}

func (r *Role) Get(ctx context.Context, id string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	c := entity.GetRoleCollection(ctx, r.Client)
	filter := DefaultFilter(ctx, Filter("_id", id))
	var item entity.Role
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaRole(), nil
}

func (r *Role) Create(ctx context.Context, item schema.Role) error {
	eitem := entity.SchemaRole(item).ToRole()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetRoleCollection(ctx, r.Client)
	err := Insert(ctx, c, eitem)
	return errors.WithStack(err)
}

func (r *Role) Update(ctx context.Context, id string, item schema.Role) error {
	eitem := entity.SchemaRole(item).ToRole()
	eitem.UpdatedAt = time.Now()
	c := entity.GetRoleCollection(ctx, r.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", id)), eitem)
	return errors.WithStack(err)
}

// Delete role
func (r *Role) Delete(ctx context.Context, id string) error {
	c := entity.GetRoleCollection(ctx, r.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", id)))
	return errors.WithStack(err)
}

// Update status role
func (r *Role) UpdateStatus(ctx context.Context, id string, status int) error {
	c := entity.GetRoleCollection(ctx, r.Client)
	err := UpdateFields(ctx, c, DefaultFilter(ctx, Filter("_id", id)), bson.M{"status": status})
	return errors.WithStack(err)
}
