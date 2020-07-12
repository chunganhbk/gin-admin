package gorm

import (
	"context"
	entity "github.com/chunganhbk/gin-go/internal/app/models/gorm"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)


// UserRole
type UserRole struct {
	DB *gorm.DB
}
func NewUserRole(db *gorm.DB) *UserRole{
	return &UserRole{db}
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

	db := entity.GetUserRoleDB(ctx, a.DB)
	if v := params.UserID; v != "" {
		db = db.Where("user_id=?", v)
	}
	if v := params.UserIDs; len(v) > 0 {
		db = db.Where("user_id IN (?)", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.UserRoles
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
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
	db := entity.GetUserRoleDB(ctx, a.DB).Where("id=?", id)
	var item entity.UserRole
	ok, err := FindOne(ctx, db, &item)
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
	result := entity.GetUserRoleDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

// Update
func (a *UserRole) Update(ctx context.Context, id string, item schema.UserRole) error {
	eitem := entity.SchemaUserRole(item).ToUserRole()
	result := entity.GetUserRoleDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete
func (a *UserRole) Delete(ctx context.Context, id string) error {
	result := entity.GetUserRoleDB(ctx, a.DB).Where("id=?", id).Delete(entity.UserRole{})
	return errors.WithStack(result.Error)
}

// Delete By UserID
func (a *UserRole) DeleteByUserID(ctx context.Context, userID string) error {
	result := entity.GetUserRoleDB(ctx, a.DB).Where("user_id=?", userID).Delete(entity.UserRole{})
	return errors.WithStack(result.Error)
}
