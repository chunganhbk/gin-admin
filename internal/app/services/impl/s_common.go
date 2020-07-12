package impl

import (
	"context"

	"github.com/chunganhbk/gin-go/internal/app/icontext"
	"github.com/chunganhbk/gin-go/internal/app/repositories"
)

// TransFunc
type TransFunc func(context.Context) error

// ExecTrans
func ExecTrans(ctx context.Context, transRp repositories.ITrans, fn TransFunc) error {
	return transRp.Exec(ctx, fn)
}

// ExecTransWithLock
func ExecTransWithLock(ctx context.Context, transRp repositories.ITrans, fn TransFunc) error {
	if !icontext.FromTransLock(ctx) {
		ctx = icontext.NewTransLock(ctx)
	}
	return ExecTrans(ctx, transRp, fn)
}

// NewNoTrans
func NewNoTrans(ctx context.Context) context.Context {
	if !icontext.FromNoTrans(ctx) {
		return icontext.NewNoTrans(ctx)
	}
	return ctx
}
