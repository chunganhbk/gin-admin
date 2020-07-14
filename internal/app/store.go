package app

import "C"
import (
	"errors"
	"github.com/chunganhbk/gin-go/internal/app/config"
	"github.com/chunganhbk/gin-go/internal/app/models"
	igorm "github.com/chunganhbk/gin-go/internal/app/repositories/gorm"
	imongo "github.com/chunganhbk/gin-go/internal/app/repositories/mongo"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
)

func InitStore(container *dig.Container) (func(), error) {
	cfg := config.C
	switch cfg.Store {
	case "gorm":
		db, storeCall, err := models.InitGormDB()
		if err != nil {
			return nil, err
		}
		_ = container.Provide(func() *gorm.DB {
			return db
		})
		_ = igorm.Inject(container)
		return storeCall, nil
	case "mongo":
		db, storeCall, err := models.InitMongo()
		if err != nil {
			return nil, err
		}

		_ = container.Provide(func() *mongo.Client {
			return db
		})
		_ = imongo.Inject(container)
		return storeCall, nil
	default:
		return nil, errors.New("unknown store")
	}

}

