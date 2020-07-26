package gorm

import (
	"time"

	"github.com/chunganhbk/gin-go/pkg/logger"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var tableName string

// Config
type Config struct {
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TableName    string
}

// New gorm
func New(c *Config) *Hook {
	tableName = c.TableName

	db, err := gorm.Open(c.DBType, c.DSN)
	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)

	db.AutoMigrate(new(LogItem))
	return &Hook{
		db: db,
	}
}

// Hook gorm
type Hook struct {
	db *gorm.DB
}

// Exec
func (h *Hook) Exec(entry *logrus.Entry) error {
	item := &LogItem{
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

	result := h.db.Create(item)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

// Close
func (h *Hook) Close() error {
	return h.db.Close()
}

// LogItem
type LogItem struct {
	ID           uint      `gorm:"column:id;primary_key;auto_increment;"` // id
	Level        string    `gorm:"column:level;size:20;index;"`
	Message      string    `gorm:"column:message;size:1024;"`
	TraceID      string    `gorm:"column:trace_id;size:128;index;"`
	UserID       string    `gorm:"column:user_id;size:36;index;"`
	SpanTitle    string    `gorm:"column:span_title;size:256;"`
	SpanFunction string    `gorm:"column:span_function;size:256;"`
	Data         string    `gorm:"column:data;type:text;"`                // (json)
	Version      string    `gorm:"column:version;index;size:32;"`
	CreatedAt    time.Time `gorm:"column:created_at;index"`
}

// TableName
func (LogItem) TableName() string {
	return tableName
}
