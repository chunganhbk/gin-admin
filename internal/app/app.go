package app

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/config"
	"github.com/chunganhbk/gin-go/internal/app/iutil"
	"github.com/chunganhbk/gin-go/internal/app/services/impl"
	"github.com/chunganhbk/gin-go/pkg/jwt"
	"github.com/chunganhbk/gin-go/pkg/logger"
	"go.uber.org/dig"
	"os"
)

type options struct {
	ConfigFile string
	ModelFile  string
	MenuFile   string
	WWWDir     string
	Version    string
}

// Option
type Option func(*options)

// SetConfigFile
func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

// SetModelFile casbin
func SetModelFile(s string) Option {
	return func(o *options) {
		o.ModelFile = s
	}
}

// SetWWWDir
func SetWWWDir(s string) Option {
	return func(o *options) {
		o.WWWDir = s
	}
}

// SetMenuFile
func SetMenuFile(s string) Option {
	return func(o *options) {
		o.MenuFile = s
	}
}

// SetVersion
func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

// Init
func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	config.MustLoad(o.ConfigFile)
	if v := o.ModelFile; v != "" {
		config.C.Casbin.Model = v
	}
	if v := o.WWWDir; v != "" {
		config.C.WWW = v
	}
	if v := o.MenuFile; v != "" {
		config.C.Menu.Data = v
	}
	config.PrintWithJSON()

	logger.Printf(ctx, "Service started, running mode：%s，version number：%s，process number：%d", config.C.RunMode, o.Version, os.Getpid())

	// Initialize unique id
	iutil.InitID()


	loggerCleanFunc, err := InitLogger()
	if err != nil {
		return nil, err
	}

	container, containerCall := BuildContainer()


	// init data
	InitData(ctx, container)


	httpServerCleanFunc := InitHTTPServer(ctx, container)

	return func() {
		httpServerCleanFunc()
		containerCall()
		loggerCleanFunc()
	}, nil
}


func BuildContainer() (*dig.Container, func()) {
	container := dig.New()

	auther, err := InitAuth()
	handleError(err)
	_ = container.Provide(func() jwt.IJWTAuth {
		return auther
	})

	//casbin
	_ = container.Provide(NewCasbinEnforcer)

	//store DB
	storeCall, err := InitStore(container)
	handleError(err)

	//register service
	err = impl.Inject(container)
	handleError(err)


	return container, func() {

		ReleaseCasbinEnforcer(container)

		if storeCall != nil {
			storeCall()
		}
	}
}
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
