package gorm

import (
	"context"
	entity "github.com/chunganhbk/gin-go/internal/app/models/gorm"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// User
type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{db}
}

func (a *User) getQueryOption(opts ...schema.UserQueryOptions) schema.UserQueryOptions {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetUserDB(ctx, a.DB)
	if v := params.Email; v != "" {
		db = db.Where("email=?", v)
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		subQuery := entity.GetUserRoleDB(ctx, a.DB).
			Select("user_id").
			Where("deleted_at is null").
			Where("role_id IN (?)", v).
			SubQuery()
		db = db.Where("id IN ?", subQuery)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("full_name LIKE ? OR phone LIKE ? OR email LIKE ?", v, v, v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.Users
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.UserQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUsers(),
	}
	return qr, nil
}

// Get
func (a *User) Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	var item entity.User
	ok, err := FindOne(ctx, entity.GetUserDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaUser(), nil
}

// Create
func (a *User) Create(ctx context.Context, item schema.User) error {
	sitem := entity.SchemaUser(item)
	result := entity.GetUserDB(ctx, a.DB).Create(sitem.ToUser())
	return errors.WithStack(result.Error)
}

// Update
func (a *User) Update(ctx context.Context, id string, item schema.User) error {
	eItem := entity.SchemaUser(item).ToUser()
	result := entity.GetUserDB(ctx, a.DB).Where("id=?", id).Updates(eItem)
	return errors.WithStack(result.Error)
}

// Delete
func (a *User) Delete(ctx context.Context, id string) error {
	result := entity.GetUserDB(ctx, a.DB).Where("id=?", id).Delete(entity.User{})
	return errors.WithStack(result.Error)
}

// Update Status
func (a *User) UpdateStatus(ctx context.Context, id string, status int) error {
	result := entity.GetUserDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}

// Update Password
func (a *User) UpdatePassword(ctx context.Context, id, password string) error {
	result := entity.GetUserDB(ctx, a.DB).Where("id=?", id).Update("password", password)
	return errors.WithStack(result.Error)
}
