package jwt

import (
	"github.com/chunganhbk/gin-go/pkg/errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type IJWTAuth interface {
	GenerateToken(userID string) (TokenInfo, error)
	ParseUserID(accessToken string, refresh bool) (string, error)
	RefreshToken(refreshToken string) (TokenInfo, error)
}

const defaultKey = "gin-go"
const defaultRefershKey = "refresh-gin-go"

var defaultOptions = options{
	tokenType:     "Bearer",
	expired:       7200,
	signingMethod: jwt.SigningMethodHS512,
	signingKey:    []byte(defaultKey),
	keyfunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrTokenInvalid
		}
		return []byte(defaultKey), nil
	},
	keyfuncRefresh: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrTokenInvalid
		}
		return []byte(defaultRefershKey), nil
	},
	expiredRefresh:    24,
	signingRefreshKey: []byte(defaultKey),
}

func NewJWTAuth(opts ...Option) *JWTAuth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	return &JWTAuth{
		opts: &o,
	}
}

type JWTAuth struct {
	opts *options
}
type options struct {
	signingMethod     jwt.SigningMethod
	signingKey        interface{}
	keyfunc           jwt.Keyfunc
	expired           int
	tokenType         string
	keyfuncRefresh    jwt.Keyfunc
	expiredRefresh    int
	signingRefreshKey interface{}
}
type Option func(*options)

func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}
func SetKeyfunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyfunc = keyFunc
	}
}

func SetSigningKey(key interface{}) Option {
	return func(o *options) {
		o.signingKey = key
	}
}
func SetExpiredRefresh(expired int) Option {
	return func(o *options) {
		o.expiredRefresh = expired
	}
}
func SetKeyfuncRefresh(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyfuncRefresh = keyFunc
	}
}

func SetSigningKeyRefresh(key interface{}) Option {
	return func(o *options) {
		o.signingRefreshKey = key
	}
}
func (jwtAuth *JWTAuth) GenerateToken(userID string) (TokenInfo, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(jwtAuth.opts.expired) * time.Second).Unix()

	token := jwt.NewWithClaims(jwtAuth.opts.signingMethod, jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt,
		NotBefore: now.Unix(),
		Subject:   userID,
	})
	tokenString, err := token.SignedString(jwtAuth.opts.signingKey)
	if err != nil {
		return nil, err
	}
	refreshToken := jwt.NewWithClaims(jwtAuth.opts.signingMethod, jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Duration(jwtAuth.opts.expired) * time.Hour).Unix(),
		NotBefore: now.Unix(),
		Subject:   userID,
	})
	rt, err := refreshToken.SignedString(jwtAuth.opts.signingRefreshKey)
	if err != nil {
		return nil, err
	}
	tokenInfo := &tokenInfo{
		ExpiresAt:    expiresAt,
		TokenType:    jwtAuth.opts.tokenType,
		AccessToken:  tokenString,
		RefreshToken: rt,
	}
	return tokenInfo, nil
}
func (a *JWTAuth) parseToken(tokenString string, refresh bool) (*jwt.StandardClaims, error) {
	option := a.opts.keyfunc
	if refresh == true {
		option = a.opts.keyfuncRefresh
	}
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, option)

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.ErrTokenMalforaled
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errors.ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.ErrTokenInvalid
			} else {
				return nil, errors.ErrTokenInvalid
			}
		}
	} else if !token.Valid {
		return nil, errors.ErrTokenInvalid
	}

	return token.Claims.(*jwt.StandardClaims), nil
}
func (jwt *JWTAuth) ParseUserID(tokenString string, refresh bool) (string, error) {
	claims, err := jwt.parseToken(tokenString, refresh)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}
func (jwtAuth *JWTAuth) RefreshToken(tokenString string) (TokenInfo, error) {
	userId, err := jwtAuth.ParseUserID(tokenString, true)
	if err != nil {
		return nil, err
	}
	return jwtAuth.GenerateToken(userId)
}
