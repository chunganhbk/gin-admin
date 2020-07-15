package app

import (
	"github.com/chunganhbk/gin-go/internal/app/adapter"
	"github.com/chunganhbk/gin-go/internal/app/repositories"
	"go.uber.org/dig"
	"time"
	"github.com/casbin/casbin/v2"
	"github.com/chunganhbk/gin-go/internal/app/config"
)
func NewCasbinEnforcer() *casbin.SyncedEnforcer {
	cfg := config.C.Casbin
	if !cfg.Enable {
		return nil
	}

	e, err := casbin.NewSyncedEnforcer(cfg.Model)
	handleError(err)

	e.EnableAutoSave(false)
	e.EnableAutoBuildRoleLinks(true)

	if cfg.Debug {
		e.EnableLog(true)
	}
	return e
}
func InitCasbinEnforcer(container *dig.Container) error {
	cfg := config.C.Casbin
	if !cfg.Enable {
		return nil
	}

	return container.Invoke(func(e *casbin.SyncedEnforcer, roleRp repositories.IRole, roleMenuRp repositories.IRoleMenu,
		menuResourceRp repositories.IMenuActionResource,
		userRp repositories.IUser, userRoleRp repositories.IUserRole,
	) error {
		adapter := adapter.NewCasbinAdapter(roleRp, roleMenuRp , menuResourceRp, userRp, userRoleRp )

		if cfg.AutoLoad {
			_ = e.InitWithModelAndAdapter(e.GetModel(), adapter)
			e.StartAutoLoadPolicy(time.Duration(cfg.AutoLoadInternal) * time.Second)
		} else {
			err := adapter.LoadPolicy(e.GetModel())
			if err != nil {
				return err
			}
		}

		err := e.BuildRoleLinks()
		if err != nil {
			return err
		}

		return nil
	})
}
func ReleaseCasbinEnforcer(container *dig.Container) {
	cfg := config.C.Casbin
	if !cfg.Enable || !cfg.AutoLoad {
		return
	}

	_ = container.Invoke(func(e *casbin.SyncedEnforcer) {
		e.StopAutoLoadPolicy()
	})
}
