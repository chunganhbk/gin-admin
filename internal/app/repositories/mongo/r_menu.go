package mongo

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"time"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	entity "github.com/chunganhbk/gin-go/internal/app/models/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Menu
type Menu struct {
	Client *mongo.Client
}
func NewMenu(client *mongo.Client) *Menu {
	return &Menu{client}
}
func (a *Menu) getQueryOption(opts ...schema.MenuQueryOptions) schema.MenuQueryOptions {
	var opt schema.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query
func (m *Menu) Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error) {
	opt := m.getQueryOption(opts...)

	c := entity.GetMenuCollection(ctx, m.Client)
	filter := DefaultFilter(ctx)
	if v := params.IDs; len(v) > 0 {
		filter = append(filter, Filter("_id", bson.M{"$in": v}))
	}
	if v := params.Name; v != "" {
		filter = append(filter, Filter("name", v))
	}
	if v := params.QueryValue; v != "" {
		filter = append(filter, Filter("$or", bson.A{
			OrRegexFilter("name", v),
			OrRegexFilter("memo", v),
		}))
	}
	if v := params.ParentID; v != nil {
		filter = append(filter, Filter("parent_id", *v))
	}
	if v := params.PrefixParentPath; v != "" {
		filter = append(filter, RegexFilter("parent_path", fmt.Sprintf("^%s.*", v)))
	}
	if v := params.ShowStatus; v != 0 {
		filter = append(filter, Filter("show_status", v))
	}
	if v := params.Status; v != 0 {
		filter = append(filter, Filter("status", v))
	}
	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByDESC))

	var list entity.Menus
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.MenuQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenus(),
	}

	return qr, nil
}

// Get menu
func (m *Menu) Get(ctx context.Context, id string, opts ...schema.MenuQueryOptions) (*schema.Menu, error) {
	c := entity.GetMenuCollection(ctx, m.Client)
	filter := DefaultFilter(ctx, Filter("_id", id))
	var item entity.Menu
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenu(), nil
}

// Create menu
func (m *Menu) Create(ctx context.Context, item schema.Menu) error {
	eitem := entity.SchemaMenu(item).ToMenu()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetMenuCollection(ctx, m.Client)
	err := Insert(ctx, c, eitem)
	return errors.WithStack(err)
}

// Update menu
func (m *Menu) Update(ctx context.Context, id string, item schema.Menu) error {
	eitem := entity.SchemaMenu(item).ToMenu()
	eitem.UpdatedAt = time.Now()
	c := entity.GetMenuCollection(ctx, m.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", id)), eitem)
	return errors.WithStack(err)
}

// Delete menu
func (m *Menu) Delete(ctx context.Context, id string) error {
	c := entity.GetMenuCollection(ctx, m.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", id)))
	return errors.WithStack(err)
}

// Update Status menu
func (m *Menu) UpdateStatus(ctx context.Context, id string, status int) error {
	c := entity.GetMenuCollection(ctx, m.Client)
	err := UpdateFields(ctx, c, DefaultFilter(ctx, Filter("_id", id)), bson.M{"status": status})
	return errors.WithStack(err)
}

// Update Parent Path
func (a *Menu) UpdateParentPath(ctx context.Context, id, parentPath string) error {
	c := entity.GetMenuCollection(ctx, a.Client)
	err := UpdateFields(ctx, c, DefaultFilter(ctx, Filter("_id", id)), bson.M{"parent_path": parentPath})
	return errors.WithStack(err)
}
