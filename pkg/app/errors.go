package app

import (
	"github.com/pkg/errors"
)

// 定义别名
var (
	New          = errors.New
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)

// 定义错误
var (
	ErrBadRequest              = New400Response("Request error")
	ErrInvalidParent           = New400Response("Invalid parent node")
	ErrNotAllowDeleteWithChild = New400Response("Contains children, cannot be deleted")
	ErrNotAllowDelete          = New400Response("Resources are not allowed to be deleted")
	ErrInvalidUserName         = New400Response("Invalid username")
	ErrInvalidPassword         = New400Response("无效的密码")
	ErrInvalidUser             = New400Response("无效的用户")
	ErrUserDisable             = New400Response("用户被禁用，请联系管理员")

	ErrNoPerm          = NewResponse(401, 401, "无访问权限")
	ErrInvalidToken    = NewResponse(9999, 401, "令牌失效")
	ErrNotFound        = NewResponse(404, 404, "资源不存在")
	ErrMethodNotAllow  = NewResponse(405, 405, "方法不被允许")
	ErrTooManyRequests = NewResponse(429, 429, "请求过于频繁")
	ErrInternalServer  = NewResponse(500, 500, "服务器发生错误")
)
