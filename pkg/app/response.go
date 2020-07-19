package app

import (
	"net/http"

	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/errors"
	"github.com/chunganhbk/gin-go/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
	ERR  error       `json:"err,omitempty"`
}

var OK = "Success"

func ResError(c *gin.Context, err error, status ...int) {
	ctx := c.Request.Context()

	var res *errors.ResponseError
	if err != nil {
		if e, ok := err.(*errors.ResponseError); ok {
			res = e
		} else {
			res = errors.UnWrapResponse(errors.Wrap500Response(err, errors.ERROR))
		}
	} else {
		res = errors.UnWrapResponse(errors.ErrInternalServer)
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
	ResJSON(c, res.StatusCode, Response{Code: res.Code, Msg: res.Message, ERR: res.ERR})
}
func ResPage(c *gin.Context, v interface{}, pr *schema.PaginationResult) {
	list := schema.ListResult{
		List:       v,
		Pagination: pr,
	}
	ResSuccess(c, list)

}
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK,
		Response{
			Code: http.StatusOK,
			Msg:  OK,
			Data: v,
		})
}
func ResList(c *gin.Context, v interface{}) {
	ResSuccess(c, schema.ListResult{List: v})
}
func ResOK(c *gin.Context) {
	ResSuccess(c, nil)
}
func ResJSON(c *gin.Context, httpCode int, res Response) {
	c.JSON(httpCode, res)
	c.AbortWithStatus(httpCode)
	return
}
