package mongo

import (
	"context"
	"time"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type TransFunc func(context.Context) error

// ExecTrans
func ExecTrans(ctx context.Context, cli *mongo.Client, fn TransFunc) error {
	transModel := &Trans{Client: cli}
	return transModel.Exec(ctx, fn)
}

// WrapPageQuery
func WrapPageQuery(ctx context.Context, c *mongo.Collection, pp schema.PaginationParam, filter interface{}, out interface{}, opts ...*options.FindOptions) (*schema.PaginationResult, error) {
	if pp.OnlyCount {
		count, err := c.CountDocuments(ctx, filter)
		if err != nil {
			return nil, err
		}
		return &schema.PaginationResult{Total: int(count)}, nil
	} else if !pp.Pagination {
		cursor, err := c.Find(ctx, filter, opts...)
		if err != nil {
			return nil, err
		}
		err = cursor.All(ctx, out)
		return nil, err
	}

	total, err := FindPage(ctx, c, pp, filter, out, opts...)
	if err != nil {
		return nil, err
	}

	return &schema.PaginationResult{
		Total:    total,
		Current:  pp.GetCurrent(),
		PageSize: pp.GetPageSize(),
	}, nil
}

// FindPage
func FindPage(ctx context.Context, c *mongo.Collection, pp schema.PaginationParam, filter interface{}, out interface{}, opts ...*options.FindOptions) (int, error) {
	count, err := c.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	} else if count == 0 {
		return 0, nil
	}

	current, pageSize := pp.GetCurrent(), pp.GetPageSize()
	opt := new(options.FindOptions)
	if len(opts) > 0 {
		opt = opts[0]
	}

	if current > 0 && pageSize > 0 {
		opt.SetSkip(int64((current - 1) * pageSize))
		opt.SetLimit(int64(pageSize))
	} else if pageSize > 0 {
		opt.SetLimit(int64(pageSize))
	}

	cursor, err := c.Find(ctx, filter, opt)
	if err != nil {
		return 0, err
	}
	err = cursor.All(ctx, out)
	return int(count), err
}

// FindOne
func FindOne(ctx context.Context, c *mongo.Collection, filter, out interface{}) (bool, error) {
	result := c.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNilDocument {
			return false, nil
		}
		return false, err
	}
	err := result.Decode(out)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Insert
func Insert(ctx context.Context, c *mongo.Collection, doc interface{}) error {
	_, err := c.InsertOne(ctx, doc)
	return err
}

// Insert Many
func InsertMany(ctx context.Context, c *mongo.Collection, docs ...interface{}) error {
	_, err := c.InsertMany(ctx, docs)
	return err
}

// Update Fields
func UpdateFields(ctx context.Context, c *mongo.Collection, filter, doc interface{}) error {
	return Update(ctx, c, filter, doc)
}

// Update Many Fields
func UpdateManyFields(ctx context.Context, c *mongo.Collection, filter, doc interface{}) error {
	return UpdateMany(ctx, c, filter, doc)
}

// Update
func Update(ctx context.Context, c *mongo.Collection, filter, doc interface{}) error {
	_, err := c.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: doc}})
	return err
}

// Update Many
func UpdateMany(ctx context.Context, c *mongo.Collection, filter, doc interface{}) error {
	_, err := c.UpdateMany(ctx, filter, bson.D{{Key: "$set", Value: doc}})
	return err
}

// Delete
func Delete(ctx context.Context, c *mongo.Collection, filter interface{}) error {
	_, err := c.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.M{"deleted_at": time.Now()}}})
	return err
}

// Delete Many
func DeleteMany(ctx context.Context, c *mongo.Collection, filter interface{}) error {
	_, err := c.UpdateMany(ctx, filter, bson.D{{Key: "$set", Value: bson.M{"deleted_at": time.Now()}}})
	return err
}

// Default Filter
func DefaultFilter(ctx context.Context, params ...bson.E) bson.D {
	var d bson.D
	if len(params) > 0 {
		d = append(d, params...)
	}
	d = append(d, Filter("deleted_at", bson.M{"$exists": 0}))
	return d
}

// RegexFilter
func RegexFilter(key, value string) bson.E {
	return bson.E{
		Key: key,
		Value: bson.M{
			"$regex":   value,
			"$options": "i",
		},
	}
}

// OrRegexFilter 正则过滤($or)
func OrRegexFilter(key, value string) bson.M {
	return bson.M{
		key: bson.M{
			"$regex":   value,
			"$options": "i",
		},
	}
}

// Filter 过滤
func Filter(key string, value interface{}) bson.E {
	return bson.E{
		Key:   key,
		Value: value,
	}
}

// Order Field Func
type OrderFieldFunc func(string) string

// Parse Order
func ParseOrder(items []*schema.OrderField, handle ...OrderFieldFunc) bson.D {
	d := make(bson.D, 0)
	for _, item := range items {
		key := item.Key
		if len(handle) > 0 {
			key = handle[0](key)
		}

		direction := 1
		if item.Direction == schema.OrderByDESC {
			direction = -1
		}
		d = append(d, bson.E{Key: key, Value: direction})
	}

	return d
}
