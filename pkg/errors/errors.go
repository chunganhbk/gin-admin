package errors

import (
	"github.com/pkg/errors"
)

const (
	SUCCESS                        = 200
	ERROR                          = 500
	INVALID_PARAMS                 = 400
	ERROR_BAD_REQUEST              = 421
	ERROR_NO_PERRMISSION           = 403
	ERROR_NOT_FOUND                = 404
	ERROR_METHOD_NOT_ALLOW         = 405
	ERROR_INVALID_PARENT           = 409
	ERROR_ALLOW_DELETE_WITH_CHILD  = 410
	ERROR_NOT_ALLOW_DELETE         = 411
	ERROR_USER_DISABLED            = 412
	ERROR_EXIST_MENU_NAME          = 413
	ERROR_EXIST_ROLE               = 414
	ERROR_EXIST_ROLE_USER          = 415
	ERROR_NOT_EXIST_USER           = 416
	ERROR_LOGIN_FAILED             = 422
	ERROR_INVALID_OLD_PASS         = 423
	ERROR_PASSWORD_REQUIRED        = 424
	ERROR_TOO_MANY_REQUEST         = 429
	ERROR_INTERNAL_SERVER          = 512
	ERROR_AUTH_CHECK_TOKEN_FAIL    = 401
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 402
	ERROR_AUTH_TOKEN               = 408
	ERROR_AUTH                     = 407
	ERROR_EXIST_EMAIL              = 430
	ERROR_NOT_EXIST_ROLE           = 431
	ERROR_TOKEN_EXPIRED            = 461
	ERROR_TOKEN_INVALID            = 462
	ERROR_TOKEN_MALFORMALED        = 463
)

var (
	New                = errors.New
	Wrap               = errors.Wrap
	Wrapf              = errors.Wrapf
	WithStack          = errors.WithStack
	WithMessage        = errors.WithMessage
	WithMessagef       = errors.WithMessagef
	ErrInternalServer  = NewResponse(ERROR_INTERNAL_SERVER, ERROR, GetMsg(ERROR_INTERNAL_SERVER))
	ErrMethodNotAllow  = NewResponse(ERROR_METHOD_NOT_ALLOW, ERROR_METHOD_NOT_ALLOW, GetMsg(ERROR_METHOD_NOT_ALLOW))
	ErrNoPermission    = NewResponse(ERROR_NO_PERRMISSION, ERROR_NO_PERRMISSION, GetMsg(ERROR_NO_PERRMISSION))
	ErrNotFound        = NewResponse(ERROR_NOT_FOUND, ERROR_NOT_FOUND, GetMsg(ERROR_NOT_FOUND))
	ErrTokenExpired    = NewStatusUnauthorized(ERROR_TOKEN_EXPIRED)
	ErrTokenInvalid    = NewStatusUnauthorized(ERROR_TOKEN_INVALID)
	ErrTokenMalforaled = NewStatusUnauthorized(ERROR_TOKEN_MALFORMALED)
)
