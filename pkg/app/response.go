package app

import (
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	ERR  error       `json:"err"`
}
const (
	SUCCESS                        = 200
	ERROR                          = 500
	INVALID_PARAMS                 = 400
	ERROR_BAD_REQUEST              = 401
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
	ERROR_TOO_MANY_REQUEST         = 429
	ERROR_INTERNAL_SERVER          = 512
	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002
	ERROR_AUTH_TOKEN               = 20003
	ERROR_AUTH                     = 20004
	ERROR_EXIST_EMAIL              = 30001
)

func ResError(c *gin.Context, err error, status ...int) {
	ctx := c.Request.Context()
	var res *ResponseError
	if err != nil {
		if e, ok := err.(*ResponseError); ok {
			res = e
		} else {
			res.Code = ERROR
			res.StatusCode = ERROR
		}
	} else {
		res.Code = ERROR_INTERNAL_SERVER
		res.StatusCode = ERROR
	}

	if len(status) > 0 {
		res.StatusCode = status[0]
	}

	if err := res.ERR; err != nil {
		if status := res.StatusCode; status >= 400 && status < 500 {
			logger.StartSpan(ctx).Warnf(err.Error())
		} else if status >= 500 {
			logger.ErrorStack(ctx, err)
		}
	}
	ResJSON(c, res.StatusCode,  res.Code, nil, err)
}
func ResPage(c *gin.Context, v interface{}, pr *schema.PaginationResult) {
	list := schema.ListResult{
		List:       v,
		Pagination: pr,
	}
	ResSuccess(c, list)

}
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, SUCCESS, v, nil)
}

func ResJSON(c *gin.Context, httpCode, errCode int, data interface{}, err error) {
	c.JSON(httpCode, Response{
		Code: errCode,
		Msg:  GetMsg(errCode),
		Data: data,
		ERR:  err,
	})
	return
}
