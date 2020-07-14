package impl

import (
	"github.com/chunganhbk/gin-go/internal/app/services"
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(NewAuthService)
	_ = container.Provide(func(m *AuthService)  services.IAuthService { return m })
	_ = container.Provide(NewMenuService)
	_ = container.Provide(func(m *MenuService)  services.IMenuService { return m })
	_ = container.Provide(NewRoleService)
	_ = container.Provide(func(m *RoleService)  services.IRoleService { return m })
	_ = container.Provide(NewUserService)
	_ = container.Provide(func(m *UserService)  services.IUserService { return m })
	return nil
}
