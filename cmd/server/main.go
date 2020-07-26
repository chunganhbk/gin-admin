package main

import (
	"context"
	"os"

	"github.com/chunganhbk/gin-go/internal/app"
	"github.com/chunganhbk/gin-go/pkg/logger"
	"github.com/urfave/cli/v2"
)

// VERSION：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "1.1.0"

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
		Usage: "Run web service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "Configuration file(.json,.yaml,.toml)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "models",
				Aliases:  []string{"m"},
				Usage:    "casbin's access control model(.conf)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "menu",
				Usage: "Initialize the menu data configuration file (.yaml)",
			},
			&cli.StringFlag{
				Name:  "www",
				Usage: "Static site directory",
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
