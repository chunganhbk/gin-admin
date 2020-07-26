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
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
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
// Run server
func Run(ctx context.Context, opts ...Option) error {
	var state int32 = 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := Init(ctx, opts...)
	if err != nil {
		return err
	}

	EXIT:
		for {
			sig := <-sc
			logger.Printf(ctx, "Received a signal[%s]", sig.String())
			switch sig {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				atomic.CompareAndSwapInt32(&state, 1, 0)
				break EXIT
			case syscall.SIGHUP:
			default:
				break EXIT
			}
		}

		cleanFunc()
		logger.Printf(ctx, "Service exit")
		time.Sleep(time.Second)
		os.Exit(int(atomic.LoadInt32(&state)))
		return nil
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
