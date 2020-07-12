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

// Menu Action Resource
type MenuActionResource struct {
	Client *mongo.Client
}
 func NewMenuActionResource (client *mongo.Client) *MenuActionResource{
	 return &MenuActionResource{client}
 }
func (a *MenuActionResource) getQueryOption(opts ...schema.MenuActionResourceQueryOptions) schema.MenuActionResourceQueryOptions {
	var opt schema.MenuActionResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query Menu Action Resource
func (a *MenuActionResource) Query(ctx context.Context, params schema.MenuActionResourceQueryParam, opts ...schema.MenuActionResourceQueryOptions) (*schema.MenuActionResourceQueryResult, error) {
	opt := a.getQueryOption(opts...)

	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	filter := DefaultFilter(ctx)
	menuIDs := params.MenuIDs
	if v := params.MenuID; v != "" {
		menuIDs = append(menuIDs, v)
	}
	if v := menuIDs; len(v) > 0 {
		actionIDs, err := a.queryActionIDs(ctx, v...)
		if err != nil {
			return nil, err
		}
		filter = append(filter, Filter("action_id", bson.M{"$in": actionIDs}))
	}
	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByASC))

	var list entity.MenuActionResources
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.MenuActionResourceQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenuActionResources(),
	}

	return qr, nil
}

// Get Menu Action Resource
func (a *MenuActionResource) Get(ctx context.Context, id string, opts ...schema.MenuActionResourceQueryOptions) (*schema.MenuActionResource, error) {
	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	filter := DefaultFilter(ctx, Filter("_id", id))
	var item entity.MenuActionResource
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenuActionResource(), nil
}

// Create Menu Action Resource
func (a *MenuActionResource) Create(ctx context.Context, item schema.MenuActionResource) error {
	eitem := entity.SchemaMenuActionResource(item).ToMenuActionResource()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	return errors.WithStack(err)
}

// Update Menu Action Resource
func (a *MenuActionResource) Update(ctx context.Context, id string, item schema.MenuActionResource) error {
	eitem := entity.SchemaMenuActionResource(item).ToMenuActionResource()
	eitem.UpdatedAt = time.Now()
	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", id)), eitem)
	return errors.WithStack(err)
}

// Delete Menu Action Resource
func (a *MenuActionResource) Delete(ctx context.Context, id string) error {
	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", id)))
	return errors.WithStack(err)
}

// Delete By Action ID
func (a *MenuActionResource) DeleteByActionID(ctx context.Context, actionID string) error {
	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	err := DeleteMany(ctx, c, DefaultFilter(ctx, Filter("action_id", actionID)))
	return errors.WithStack(err)
}

// Delete By MenuID
func (a *MenuActionResource) DeleteByMenuID(ctx context.Context, menuID string) error {
	actionIDs, err := a.queryActionIDs(ctx, menuID)
	if err != nil {
		return err
	}

	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	err = DeleteMany(ctx, c, DefaultFilter(ctx, Filter("action_id", bson.M{"$in": actionIDs})))
	return errors.WithStack(err)
}

func (a *MenuActionResource) queryActionIDs(ctx context.Context, menuIDs ...string) ([]interface{}, error) {
	result, err := entity.GetMenuActionCollection(ctx, a.Client).Distinct(ctx, "_id", DefaultFilter(ctx, Filter("menu_id", bson.M{"$in": menuIDs})))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}
