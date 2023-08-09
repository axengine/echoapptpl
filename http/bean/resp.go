package bean

import (
	"echoapptpl/types/errorx"
	"github.com/axengine/utils/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Resp http resp
type Resp struct {
	ResCode int         `json:"resCode"`
	ResDesc string      `json:"resDesc"`
	Result  interface{} `json:"result"`
	TraceId string      `json:"traceId,omitempty"`
}

// ResultPage result with page info
type ResultPage struct {
	Content interface{} `json:"content"`
	Total   int64       `json:"total"`
}

// Success display successful signal
func (r *Resp) Success(result interface{}) *Resp {
	r.ResCode = 0
	r.ResDesc = "ok"
	r.Result = result
	return r
}

// Success display successful result with page info
func (r *Resp) SuccessPage(content interface{}, total int64) *Resp {
	r.ResCode = 0
	r.ResDesc = "ok"
	r.Result = ResultPage{content, total}
	return r
}

// Fail display fail signal
func (r *Resp) Fail(code int, desc string, result interface{}) *Resp {
	r.ResCode = code
	r.ResDesc = desc
	r.Result = result
	return r
}

// FailMsg display fail msg
func (r *Resp) FailMsg(desc string) *Resp {
	return r.Fail(400, desc, nil)
}

func (r *Resp) FailErr(c echo.Context, err error) *Resp {
	var (
		code int
		msg  string
	)
	r.TraceId = uuid.New().String()
	log.Logger.Warn("FailErr", zap.String("traceId", r.TraceId), zap.Error(err))
	if e, ok := err.(errorx.Error); ok {
		code = e.Code
		//msg = errorx.GetMessage(e.Code)
		msg = e.String() // 只返回错误说明给前端
	} else {
		code = errorx.CodeSystem
		msg = "system error"
	}

	r.ResCode = code
	r.ResDesc = msg
	return r
}
