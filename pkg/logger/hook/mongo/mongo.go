package bson

import (
	"context"
	"time"

	"github.com/chunganhbk/gin-go/pkg/logger"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config
type Config struct {
	URI        string
	Database   string
	Collection string
	Timeout    time.Duration
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// New bson
func New(cfg *Config) *Hook {
	var (
		ctx    = context.Background()
		cancel context.CancelFunc
	)

	if t := cfg.Timeout; t > 0 {
		ctx, cancel = context.WithTimeout(ctx, t)
		defer cancel()
	}

	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	handleError(err)
	c := cli.Database(cfg.Database).Collection(cfg.Collection)

	return &Hook{
		Client:     cli,
		Collection: c,
	}
}

// Hook bson
type Hook struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

// Exec
func (h *Hook) Exec(entry *logrus.Entry) error {
	item := &LogItem{
		ID:        primitive.NewObjectID().Hex(),
		Level:     entry.Level.String(),
		Message:   entry.Message,
		CreatedAt: entry.Time,
	}

	data := entry.Data
	if v, ok := data[logger.TraceIDKey]; ok {
		item.TraceID, _ = v.(string)
		delete(data, logger.TraceIDKey)
	}
	if v, ok := data[logger.UserIDKey]; ok {
		item.UserID, _ = v.(string)
		delete(data, logger.UserIDKey)
	}
	if v, ok := data[logger.SpanTitleKey]; ok {
		item.SpanTitle, _ = v.(string)
		delete(data, logger.SpanTitleKey)
	}
	if v, ok := data[logger.SpanFunctionKey]; ok {
		item.SpanFunction, _ = v.(string)
		delete(data, logger.SpanFunctionKey)
	}
	if v, ok := data[logger.VersionKey]; ok {
		item.Version, _ = v.(string)
		delete(data, logger.VersionKey)
	}

	if len(data) > 0 {
		item.Data = util.JSONMarshalToString(data)
	}

	_, err := h.Collection.InsertOne(context.Background(), item)
	return err
}

// Close
func (h *Hook) Close() error {
	return h.Client.Disconnect(context.Background())
}

// LogItem
type LogItem struct {
	ID           string    `bson:"_id"`           // id
	Level        string    `bson:"level"`
	Message      string    `bson:"message"`
	TraceID      string    `bson:"trace_id"`
	UserID       string    `bson:"user_id"`
	SpanTitle    string    `bson:"span_title"`
	SpanFunction string    `bson:"span_function"`
	Data         string    `bson:"data"`          //(json)
	Version      string    `bson:"version"`
	CreatedAt    time.Time `bson:"created_at"`
}
