package gorm

import (
	"context"
	"github.com/pkg/errors"
	entity"github.com/chunganhbk/gin-go/internal/app/models/gorm"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/jinzhu/gorm"
)


// Role \
type Role struct {
	DB *gorm.DB
}
func NewRole(db *gorm.DB) *Role{
	return &Role{db}
}
func (a *Role) getQueryOption(opts ...schema.RoleQueryOptions) schema.RoleQueryOptions {
	var opt schema.RoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query
func (a *Role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetRoleDB(ctx, a.DB)
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name=?", v)
	}
	if v := params.UserID; v != "" {
		subQuery := entity.GetUserRoleDB(ctx, a.DB).
			Where("deleted_at is null").
			Where("user_id=?", v).
			Select("role_id").SubQuery()
		db = db.Where("id IN ?", subQuery)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR memo LIKE ?", v, v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.Roles
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoles(),
	}

	return qr, nil
}

// Get
func (a *Role) Get(ctx context.Context, id string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	var role entity.Role
	ok, err := FindOne(ctx, entity.GetRoleDB(ctx, a.DB).Where("id=?", id), &role)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return role.ToSchemaRole(), nil
}

// Create
func (a *Role) Create(ctx context.Context, item schema.Role) error {
	eitem := entity.SchemaRole(item).ToRole()
	result := entity.GetRoleDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

// Update
func (a *Role) Update(ctx context.Context, id string, item schema.Role) error {
	eitem := entity.SchemaRole(item).ToRole()
	result := entity.GetRoleDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete
func (a *Role) Delete(ctx context.Context, id string) error {
	result := entity.GetRoleDB(ctx, a.DB).Where("id=?", id).Delete(entity.Role{})
	return errors.WithStack(result.Error)
}

// Update Status
func (a *Role) UpdateStatus(ctx context.Context, id string, status int) error {
	result := entity.GetRoleDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
