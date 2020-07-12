package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)
type IJWTAuth interface {
	GenerateToken(userID string) (TokenInfo, error)
	ParseUserID( accessToken string) (string, error)
}

var (
	ErrInvalidToken = errors.New("invalid token")
)
const defaultKey = "gin-go"

var defaultOptions = options{
	tokenType:     "Bearer",
	expired:       7200,
	signingMethod: jwt.SigningMethodHS512,
	signingKey:    []byte(defaultKey),
	keyfunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(defaultKey), nil
	},
}
func NewJWTAuth( opts ...Option) *JWTAuth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	return &JWTAuth{
		opts:  &defaultOptions,
	}
}

type JWTAuth struct {
	opts  *options
}
type options struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyfunc       jwt.Keyfunc
	expired       int
	tokenType     string
}
type Option func(*options)
func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}
func SetSigningKey(key interface{}) Option {
	return func(o *options) {
		o.signingKey = key
	}
}

func (jwtAuth *JWTAuth) GenerateToken( userID string) (TokenInfo, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(jwtAuth.opts.expired) * time.Second).Unix()

	token := jwt.NewWithClaims(jwtAuth.opts.signingMethod, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt,
		NotBefore: now.Unix(),
		Subject:   userID,
	})

	tokenString, err := token.SignedString(jwtAuth.opts.signingKey)
	if err != nil {
		return nil, err
	}

	tokenInfo := &tokenInfo{
		ExpiresAt:   expiresAt,
		TokenType:   jwtAuth.opts.tokenType,
		AccessToken: tokenString,
	}
	return tokenInfo, nil
}
func (a *JWTAuth) parseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, a.opts.keyfunc)
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token.Claims.(*jwt.StandardClaims), nil
}
func (jwt *JWTAuth) ParseUserID( tokenString string) (string, error) {
	claims, err := jwt.parseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}
func (jwtAuth *JWTAuth) RefreshToken(tokenString string) (TokenInfo, error) {
	userId, err := jwtAuth.ParseUserID(tokenString)
	if err != nil {
		return nil, err
	}
	return jwtAuth.GenerateToken(userId)
}
