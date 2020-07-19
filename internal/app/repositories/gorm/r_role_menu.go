package gorm

import (
	"context"
	entity "github.com/chunganhbk/gin-go/internal/app/models/gorm"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// RoleMenu
type RoleMenu struct {
	DB *gorm.DB
}

func NewRoleMenu(db *gorm.DB) *RoleMenu {
	return &RoleMenu{db}
}
func (a *RoleMenu) getQueryOption(opts ...schema.RoleMenuQueryOptions) schema.RoleMenuQueryOptions {
	var opt schema.RoleMenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query
func (a *RoleMenu) Query(ctx context.Context, params schema.RoleMenuQueryParam, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenuQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetRoleMenuDB(ctx, a.DB)
	if v := params.RoleID; v != "" {
		db = db.Where("role_id=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		db = db.Where("role_id IN (?)", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.RoleMenus
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleMenuQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoleMenus(),
	}

	return qr, nil
}

// Get role menu
func (a *RoleMenu) Get(ctx context.Context, id string, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenu, error) {
	db := entity.GetRoleMenuDB(ctx, a.DB).Where("id=?", id)
	var item entity.RoleMenu
	ok, err := FindOne(ctx, db, &item)
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
	result := entity.GetRoleMenuDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

// Update
func (a *RoleMenu) Update(ctx context.Context, id string, item schema.RoleMenu) error {
	eitem := entity.SchemaRoleMenu(item).ToRoleMenu()
	result := entity.GetRoleMenuDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete
func (a *RoleMenu) Delete(ctx context.Context, id string) error {
	result := entity.GetRoleMenuDB(ctx, a.DB).Where("id=?", id).Delete(entity.RoleMenu{})
	return errors.WithStack(result.Error)
}

// Delete By RoleID
func (a *RoleMenu) DeleteByRoleID(ctx context.Context, roleID string) error {
	result := entity.GetRoleMenuDB(ctx, a.DB).Where("role_id=?", roleID).Delete(entity.RoleMenu{})
	return errors.WithStack(result.Error)
}
