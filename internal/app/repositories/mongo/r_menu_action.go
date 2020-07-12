package mongo

import (
	"context"
	"github.com/pkg/errors"
	"time"
	"github.com/chunganhbk/gin-go/internal/app/schema"
     imodel "github.com/chunganhbk/gin-go/internal/app/models/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


// Menu Action
type MenuAction struct {
	Client *mongo.Client
}
func NewMenuAction(client *mongo.Client) *MenuAction {
	return &MenuAction{client}
}
func (a *MenuAction) getQueryOption(opts ...schema.MenuActionQueryOptions) schema.MenuActionQueryOptions {
	var opt schema.MenuActionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query menu
func (a *MenuAction) Query(ctx context.Context, params schema.MenuActionQueryParam, opts ...schema.MenuActionQueryOptions) (*schema.MenuActionQueryResult, error) {
	opt := a.getQueryOption(opts...)

	c := imodel.GetMenuActionCollection(ctx, a.Client)
	filter := DefaultFilter(ctx)
	if v := params.MenuID; v != "" {
		filter = append(filter, Filter("menu_id", v))
	}
	if v := params.IDs; len(v) > 0 {
		filter = append(filter, Filter("_id", bson.M{"$in": v}))
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByASC))

	var list imodel.MenuActions
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.MenuActionQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenuActions(),
	}

	return qr, nil
}

// Get menu action
func (m *MenuAction) Get(ctx context.Context, id string, opts ...schema.MenuActionQueryOptions) (*schema.MenuAction, error) {
	c := imodel.GetMenuActionCollection(ctx, m.Client)
	filter := DefaultFilter(ctx, Filter("_id", id))
	var item imodel.MenuAction
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenuAction(), nil
}

// Create  menu action
func (m *MenuAction) Create(ctx context.Context, item schema.MenuAction) error {
	eitem := imodel.SchemaMenuAction(item).ToMenuAction()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := imodel.GetMenuActionCollection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	return errors.WithStack(err)
}

// Update menu action
func (a *MenuAction) Update(ctx context.Context, id string, item schema.MenuAction) error {
	eitem := imodel.SchemaMenuAction(item).ToMenuAction()
	eitem.UpdatedAt = time.Now()
	c := imodel.GetMenuActionCollection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", id)), eitem)
	return errors.WithStack(err)
}

// Delete menu action
func (a *MenuAction) Delete(ctx context.Context, id string) error {
	c := imodel.GetMenuActionCollection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", id)))
	return errors.WithStack(err)
}

// Delete ByMenuID
func (a *MenuAction) DeleteByMenuID(ctx context.Context, menuID string) error {
	c := imodel.GetMenuActionCollection(ctx, a.Client)
	err := DeleteMany(ctx, c, DefaultFilter(ctx, Filter("menu_id", menuID)))
	return errors.WithStack(err)
}
