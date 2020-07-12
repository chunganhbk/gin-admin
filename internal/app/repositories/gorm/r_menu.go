package gorm

import (
	"context"
	entity "github.com/chunganhbk/gin-go/internal/app/models/gorm"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)



// Menu
type Menu struct {
	DB *gorm.DB
}
func NewMenu(db *gorm.DB) *Menu{
 return &Menu{db}
}
func (a *Menu) getQueryOption(opts ...schema.MenuQueryOptions) schema.MenuQueryOptions {
	var opt schema.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query
func (a *Menu) Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetMenuDB(ctx, a.DB)
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name=?", v)
	}
	if v := params.ParentID; v != nil {
		db = db.Where("parent_id=?", *v)
	}
	if v := params.PrefixParentPath; v != "" {
		db = db.Where("parent_path LIKE ?", v+"%")
	}
	if v := params.ShowStatus; v != 0 {
		db = db.Where("show_status=?", v)
	}
	if v := params.Status; v != 0 {
		db = db.Where("status=?", v)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR memo LIKE ?", v, v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.Menus
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.MenuQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenus(),
	}

	return qr, nil
}

// Get
func (a *Menu) Get(ctx context.Context, id string, opts ...schema.MenuQueryOptions) (*schema.Menu, error) {
	var item entity.Menu
	ok, err := FindOne(ctx, entity.GetMenuDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenu(), nil
}

// Create
func (a *Menu) Create(ctx context.Context, item schema.Menu) error {
	eitem := entity.SchemaMenu(item).ToMenu()
	result := entity.GetMenuDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

// Update
func (a *Menu) Update(ctx context.Context, id string, item schema.Menu) error {
	eitem := entity.SchemaMenu(item).ToMenu()
	result := entity.GetMenuDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Update Parent Path
func (a *Menu) UpdateParentPath(ctx context.Context, id, parentPath string) error {
	result := entity.GetMenuDB(ctx, a.DB).Where("id=?", id).Update("parent_path", parentPath)
	return errors.WithStack(result.Error)
}

// Delete
func (a *Menu) Delete(ctx context.Context, id string) error {
	result := entity.GetMenuDB(ctx, a.DB).Where("id=?", id).Delete(entity.Menu{})
	return errors.WithStack(result.Error)
}

// Update Status
func (a *Menu) UpdateStatus(ctx context.Context, id string, status int) error {
	result := entity.GetMenuDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
