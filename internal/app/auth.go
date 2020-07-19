package app

import (
	"github.com/chunganhbk/gin-go/internal/app/config"
	"github.com/chunganhbk/gin-go/pkg/errors"
	jwtAuth "github.com/chunganhbk/gin-go/pkg/jwt"
	"github.com/dgrijalva/jwt-go"
)

func InitAuth() (jwtAuth.IJWTAuth, error) {

	cfg := config.C.JWTAuth
	var opts []jwtAuth.Option
	//access token
	opts = append(opts, jwtAuth.SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrTokenInvalid
		}
		return []byte(cfg.SigningKey), nil
	}))
	opts = append(opts, jwtAuth.SetExpired(cfg.Expired))
	opts = append(opts, jwtAuth.SetSigningKey([]byte(cfg.SigningKey)))
	//refresh token
	opts = append(opts, jwtAuth.SetKeyfuncRefresh(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrTokenInvalid
		}
		return []byte(cfg.SigningRefreshKey), nil
	}))
	opts = append(opts, jwtAuth.SetExpiredRefresh(cfg.ExpiredRefreshToken))
	opts = append(opts, jwtAuth.SetSigningKeyRefresh([]byte(cfg.SigningRefreshKey)))
	return jwtAuth.NewJWTAuth(opts...), nil

}
