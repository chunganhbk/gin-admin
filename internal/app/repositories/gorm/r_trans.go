package gorm

import (
	"context"
	"github.com/pkg/errors"

	"github.com/chunganhbk/gin-go/internal/app/icontext"
	"github.com/jinzhu/gorm"
)
// Trans
type Trans struct {
	DB *gorm.DB
}
func NewTrans (db *gorm.DB) *Trans{
	return &Trans{db}
}

// Exec
func (a *Trans) Exec(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := icontext.FromTrans(ctx); ok {
		return fn(ctx)
	}

	err := a.DB.Transaction(func(db *gorm.DB) error {
		return fn(icontext.NewTrans(ctx, db))
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
