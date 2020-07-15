package api

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(NewAuth)
	_ = container.Provide(NewMenu)
	_ = container.Provide(NewRole)
	_ = container.Provide(NewUser)
	return nil
}

