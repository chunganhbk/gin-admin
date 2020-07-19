package app

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/config"
	"github.com/chunganhbk/gin-go/internal/app/services"
	"go.uber.org/dig"
)

func InitData(ctx context.Context, container *dig.Container) error {
	if c := config.C.Menu; c.Enable && c.Data != "" {
		return container.Invoke(func(menu services.IMenuService, user services.IUserService) error {
			_ = menu.InitData(ctx, c.Data)
			_ = user.InitData(ctx)
			return nil
		})
	}
	return nil
}
