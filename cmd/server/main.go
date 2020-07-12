package main

import (
	"context"
	"os"

	"github.com/chunganhbk/gin-go/internal/app"
	"github.com/chunganhbk/gin-go/pkg/logger"
	"github.com/urfave/cli/v2"
)

// VERSION：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "1.0.1"

func main() {
	logger.SetVersion(VERSION)
	ctx := logger.NewTraceIDContext(context.Background(), "main")

	app := cli.NewApp()
	app.Name = "server"
	app.Version = VERSION
	app.Usage = "RBAC scaffolding based on GIN + GORM/MONGO + CASBIN + WIRE."
	app.Commands = []*cli.Command{
		newWebCmd(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Errorf(ctx, err.Error())
	}
}

func newWebCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "运行web服务",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "配置文件(.json,.yaml,.toml)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "models",
				Aliases:  []string{"m"},
				Usage:    "casbin的访问控制模型(.conf)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "menu",
				Usage: "初始化菜单数据配置文件(.yaml)",
			},
			&cli.StringFlag{
				Name:  "www",
				Usage: "静态站点目录",
			},
		},
		Action: func(c *cli.Context) error {
			return app.Run(ctx,
				app.SetConfigFile(c.String("conf")),
				app.SetModelFile(c.String("models")),
				app.SetWWWDir(c.String("www")),
				app.SetMenuFile(c.String("menu")),
				app.SetVersion(VERSION))
		},
	}
}
