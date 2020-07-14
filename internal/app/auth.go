package app

import (
	"github.com/chunganhbk/gin-go/internal/app/config"
	"github.com/chunganhbk/gin-go/pkg/jwt"
)

func InitAuth() (jwt.IJWTAuth, error) {

	jwtConfig := config.C.JWTAuth
	var opts []jwt.Option

	opts = append(opts, jwt.SetExpired(jwtConfig.Expired))
	opts = append(opts, jwt.SetSigningKey([]byte(jwtConfig.SigningKey)))
	return jwt.NewJWTAuth(opts...), nil

}
