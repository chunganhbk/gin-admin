package mongo

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/icontext"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)


// Trans
type Trans struct {
	Client *mongo.Client
}
func NewTrans (client *mongo.Client) *Trans{
	return &Trans{client}
}
// Exec
func (a *Trans) Exec(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := icontext.FromTrans(ctx); ok {
		return fn(ctx)
	}

	session, err := a.Client.StartSession()
	if err != nil {
		return errors.WithStack(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		err := fn(icontext.NewTrans(sessCtx, true))
		return nil, err
	})

	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
