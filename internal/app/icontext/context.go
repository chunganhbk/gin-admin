package icontext

import (
	"context"
)


type (
	transCtx     struct{}
	noTransCtx   struct{}
	transLockCtx struct{}
	userIDCtx    struct{}
)

// New Transaction
func NewTrans(ctx context.Context, trans interface{}) context.Context {
	return context.WithValue(ctx, transCtx{}, trans)
}

// From Trans
func FromTrans(ctx context.Context) (interface{}, bool) {
	v := ctx.Value(transCtx{})
	return v, v != nil
}

// New No Transaction
func NewNoTrans(ctx context.Context) context.Context {
	return context.WithValue(ctx, noTransCtx{}, true)
}

// From No Transaction
func FromNoTrans(ctx context.Context) bool {
	v := ctx.Value(noTransCtx{})
	return v != nil && v.(bool)
}

// New Transaction Lock
func NewTransLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, transLockCtx{}, true)
}

// From Transaction Lock
func FromTransLock(ctx context.Context) bool {
	v := ctx.Value(transLockCtx{})
	return v != nil && v.(bool)
}

// SET New UserID
func NewUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDCtx{}, userID)
}

// From UserID
func FromUserID(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s, s != ""
		}
	}
	return "", false
}
