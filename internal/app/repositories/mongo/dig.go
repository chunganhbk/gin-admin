package mongo

import (
	"github.com/chunganhbk/gin-go/internal/app/repositories"
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(NewMenu)
	_ = container.Provide(func(m *Menu)  repositories.IMenu { return m })
	_ = container.Provide(NewMenuAction)
	_ = container.Provide(func(m *MenuAction)  repositories.IMenuAction { return m })
	_ = container.Provide(NewMenuActionResource)
	_ = container.Provide(func(m *MenuActionResource)  repositories.IMenuActionResource { return m })
	_ = container.Provide(NewRole)
	_ = container.Provide(func(m *Role)  repositories.IRole { return m })
	_ = container.Provide(NewRoleMenu)
	_ = container.Provide(func(m *RoleMenu)  repositories.IRoleMenu { return m })
	_ = container.Provide(NewTrans)
	_ = container.Provide(func(m *Trans)  repositories.ITrans { return m })
	_ = container.Provide(NewUser)
	_ = container.Provide(func(m *User)  repositories.IUser { return m })
	_ = container.Provide(NewUserRole)
	_ = container.Provide(func(m *UserRole)  repositories.IUserRole { return m })
	return nil
}
