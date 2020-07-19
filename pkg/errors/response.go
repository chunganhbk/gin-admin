package errors

import (
	"fmt"
	"net/http"
)

type ResponseError struct {
	Code       int
	StatusCode int
	Message    string
	ERR        error
}

func (r *ResponseError) Error() string {
	if r.ERR != nil {
		return r.ERR.Error()
	}
	return GetMsg(r.Code)
}

func Wrap400Response(err error) error {
	return WrapResponse(err, INVALID_PARAMS, INVALID_PARAMS, err.Error())
}
func WrapUnauthorized(err error) error {
	return WrapResponse(err, http.StatusUnauthorized, ERROR_AUTH_CHECK_TOKEN_FAIL, err.Error())
}
func Wrap500Response(err error, errorCode int, args ...interface{}) error {
	return WrapResponse(err, ERROR, errorCode, args...)
}
func WrapResponse(err error, code, statusCode int, args ...interface{}) error {
	res := &ResponseError{
		Code:       code,
		Message:    fmt.Sprintf(GetMsg(code), args...),
		ERR:        err,
		StatusCode: statusCode,
	}
	return res
}

func New400Response(errCode int) error {
	return NewResponse(errCode, INVALID_PARAMS, GetMsg(errCode))
}
func New400NoErr(errCode int) error {
	return NewResponse(errCode, INVALID_PARAMS, GetMsg(errCode))
}
func NewStatusUnauthorized(code int) error {
	return NewResponse(code, http.StatusUnauthorized, GetMsg(code))
}

func NewResponse(code, statusCode int, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:       code,
		Message:    fmt.Sprintf(msg, args...),
		StatusCode: statusCode,
	}
	return res
}
func UnWrapResponse(err error) *ResponseError {
	if v, ok := err.(*ResponseError); ok {
		return v
	}
	return nil
}
