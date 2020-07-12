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

// User Role
type UserRole struct {
	Client *mongo.Client
}
func NewUserRole(client *mongo.Client) *UserRole{
	return &UserRole{client}
}
func (a *UserRole) getQueryOption(opts ...schema.UserRoleQueryOptions) schema.UserRoleQueryOptions {
	var opt schema.UserRoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query
func (a *UserRole) Query(ctx context.Context, params schema.UserRoleQueryParam, opts ...schema.UserRoleQueryOptions) (*schema.UserRoleQueryResult, error) {
	opt := a.getQueryOption(opts...)

	c := entity.GetUserRoleCollection(ctx, a.Client)
	filter := DefaultFilter(ctx)
	userIDs := params.UserIDs
	if v := params.UserID; v != "" {
		userIDs = append(userIDs, v)
	}
	if v := userIDs; len(v) > 0 {
		filter = append(filter, Filter("user_id", bson.M{"$in": v}))
	}
	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByDESC))

	var list entity.UserRoles
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.UserRoleQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUserRoles(),
	}

	return qr, nil
}

// Get
func (a *UserRole) Get(ctx context.Context, id string, opts ...schema.UserRoleQueryOptions) (*schema.UserRole, error) {
	c := entity.GetUserRoleCollection(ctx, a.Client)
	filter := DefaultFilter(ctx, Filter("_id", id))
	var item entity.UserRole
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaUserRole(), nil
}

// Create
func (a *UserRole) Create(ctx context.Context, item schema.UserRole) error {
	eitem := entity.SchemaUserRole(item).ToUserRole()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetUserRoleCollection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	return errors.WithStack(err)
}

// Update
func (a *UserRole) Update(ctx context.Context, id string, item schema.UserRole) error {
	eitem := entity.SchemaUserRole(item).ToUserRole()
	eitem.UpdatedAt = time.Now()
	c := entity.GetUserRoleCollection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", id)), eitem)
	return errors.WithStack(err)
}

// Delete
func (a *UserRole) Delete(ctx context.Context, id string) error {
	c := entity.GetUserRoleCollection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", id)))
	return errors.WithStack(err)
}

// DeleteByUserID
func (a *UserRole) DeleteByUserID(ctx context.Context, userID string) error {
	c := entity.GetUserRoleCollection(ctx, a.Client)
	err := DeleteMany(ctx, c, DefaultFilter(ctx, Filter("user_id", userID)))
	return errors.WithStack(err)
}
