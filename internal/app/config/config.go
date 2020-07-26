package config

import (
	"fmt"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/koding/multiconfig"
	"os"
	"strings"
	"sync"
)

var (
	C    = new(Config)
	once sync.Once
)

// MustLoad
func MustLoad(fpaths ...string) {
	once.Do(func() {
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
			&multiconfig.EnvironmentLoader{},
		}

		for _, fpath := range fpaths {
			if strings.HasSuffix(fpath, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "json") {
				loaders = append(loaders, &multiconfig.JSONLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "yaml") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: fpath})
			}
		}

		m := multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}
		m.MustLoad(C)
	})
}

// Print With JSON
func PrintWithJSON() {
	if C.PrintConfig {
		b, err := util.JSONMarshalIndent(C, "", " ")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}
		os.Stdout.WriteString(string(b) + "\n")
	}
}

// Config
type Config struct {
	RunMode      string
	WWW          string
	Version      string
	Swagger      bool
	PrintConfig  bool
	Store        string
	HTTP         HTTP
	Menu         Menu
	Casbin       Casbin
	Log          Log
	GrpcServer   GrpcServer
	LogGormHook  LogGormHook
	LogMongoHook LogMongoHook
	JWTAuth      JWTAuth
	Monitor      Monitor
	CORS         CORS
	GZIP         GZIP
	Gorm         Gorm
	MySQL        MySQL
	Postgres     Postgres
	Sqlite3      Sqlite3
	Mongo        Mongo
	UniqueID     struct {
		Type      string
		Snowflake struct {
			Node  int64
			Epoch int64
		}
	}
}
type GrpcServer struct {
	Enable bool
	Port int

}
// IsDebugMode
func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

// Menu
type Menu struct {
	Enable bool
	Data   string
}

// Casbin
type Casbin struct {
	Enable           bool
	Debug            bool
	Model            string
	AutoLoad         bool
	AutoLoadInternal int
}

// LogHook
type LogHook string

// IsGorm gorm
func (h LogHook) IsGorm() bool {
	return h == "gorm"
}

// IsMongo mongo
func (h LogHook) IsMongo() bool {
	return h == "mongo"
}

// Log
type Log struct {
	Level         int
	Format        string
	Output        string
	OutputFile    string
	EnableHook    bool
	HookLevels    []string
	Hook          LogHook
	HookMaxThread int
	HookMaxBuffer int
}

// Log Gorm Hook
type LogGormHook struct {
	DBType       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	Table        string
}

// Log Mongo Hook
type LogMongoHook struct {
	Collection string
}

// JWTAuth
type JWTAuth struct {
	SigningKey          string
	Expired             int //Seconds
	SigningRefreshKey   string
	ExpiredRefreshToken int //hours
}

// HTTP http
type HTTP struct {
	Host             string
	Port             int
	CertFile         string
	KeyFile          string
	ShutdownTimeout  int
	MaxContentLength int64
}

// Monitor 监控配置参数
type Monitor struct {
	Enable    bool
	Addr      string
	ConfigDir string
}

// CORS
type CORS struct {
	Enable           bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

// GZIP gzip压缩
type GZIP struct {
	Enable             bool
	ExcludedExtentions []string
	ExcludedPaths      []string
}

// Gorm gorm
type Gorm struct {
	Debug             bool
	DBType            string
	MaxLifetime       int
	MaxOpenConns      int
	MaxIdleConns      int
	TablePrefix       string
	EnableAutoMigrate bool
}

// MySQL mysql
type MySQL struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	Parameters string
}

// DSN connect mysql
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

// Postgres postgres
type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DSN connect Postgres
func (a Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		a.Host, a.Port, a.User, a.DBName, a.Password, a.SSLMode)
}

// Sqlite3
type Sqlite3 struct {
	Path string
}

// DSN connect Sqlite3
func (a Sqlite3) DSN() string {
	return a.Path
}

// Mongo mongo
type Mongo struct {
	URI              string
	Database         string
	Timeout          int
	CollectionPrefix string
}
