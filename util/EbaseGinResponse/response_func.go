package EbaseGinResponse

import (
	"github.com/gin-gonic/gin"
	"github.com/jilin7105/ebase/logger"
	"github.com/jilin7105/ebase/util/LinkTracking"
	"net/http"
)

var Default = &response{}

func Error(c *gin.Context, code int, err error, msg string) {
	EbaseRequestID := LinkTracking.GetEbaseRequestID(c)
	res := Default.Clone()
	res.SetInfo(msg)
	if err != nil {
		res.SetInfo(err.Error())
	}
	if msg != "" {
		res.SetMsg(msg)
	}
	res.SetEbaseRequestID(EbaseRequestID)
	res.SetCode(int32(code))
	res.SetSuccess(false)
	// 记录日志
	logger.Error("EbaseRequestID[%s]:errmsg(%s) ", EbaseRequestID, err.Error())
	// 写入上下文
	c.Set("result", res)
	// 返回结果集
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func OK(c *gin.Context, data any, msg string) {
	res := Default.Clone()
	EbaseRequestID := LinkTracking.GetEbaseRequestID(c)
	res.SetData(data)
	res.SetSuccess(true)
	if msg != "" {
		res.SetMsg(msg)
		res.SetInfo(msg)
	}
	res.SetEbaseRequestID(EbaseRequestID)
	res.SetCode(http.StatusOK)
	c.Set("result", res)
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func PageOK(c *gin.Context, result any, count int, pageIndex int, pageSize int, msg string) {
	var res page
	res.List = result
	res.Count = count
	res.PageIndex = pageIndex
	res.PageSize = pageSize
	OK(c, res, msg)
}
