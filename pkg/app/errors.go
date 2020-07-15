package app

import (
	"github.com/pkg/errors"
	"net/http"
)
type ResponseError struct {
	Code       int
	StatusCode int
	ERR        error
}

func (r *ResponseError) Error() string {
	if r.ERR != nil {
		return r.ERR.Error()
	}
	return GetMsg(r.Code)
}
func NewResponse(code, statusCode int, err error) error {

	res := &ResponseError{
		Code:       code,
		StatusCode: statusCode,
		ERR: err,
	}
	return res
}
func ResponseNotFound() error {

	res := &ResponseError{
		Code:     ERROR_NOT_FOUND  ,
		StatusCode: ERROR_NOT_FOUND,
	}
	return res
}
func New400Response(code int, err error) error {
	println("string %s", err.Error())
	return NewResponse(code, INVALID_PARAMS, err)
}
func NewStatusUnauthorized(code int) error {
	return NewResponse(code, http.StatusUnauthorized, nil)
}
func NoPermissionResponse() error {
	return NewResponse(ERROR_NO_PERRMISSION, ERROR_NO_PERRMISSION, nil)
}
func MethodNotAllowResponse() error{
	return NewResponse(ERROR_METHOD_NOT_ALLOW, ERROR_METHOD_NOT_ALLOW, nil)
}
var (
	New          = errors.New
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)
