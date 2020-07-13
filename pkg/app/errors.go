package app

import (
	"github.com/pkg/errors"
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
func NewResponse(code, statusCode int) error {

	res := &ResponseError{
		Code:       code,
		StatusCode: statusCode,
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
func New400Response() error {
	return NewResponse(INVALID_PARAMS, INVALID_PARAMS)
}
var (
	New          = errors.New
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)
