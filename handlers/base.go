package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Gin struct {
	*gin.Context
}

type JSONResult struct {
	Code    int         `json:"code" description:"code"`
	Message string      `json:"message,omitempty" description:"message"`
	Data    interface{} `json:"data,omitempty" description:"data"`
}

func NewRes(data interface{}, err error, code int) *JSONResult {
	if code == 200 {
		code = 0
	}
	codeMsg := getMsg(code)
	return &JSONResult{
		Data: data,
		Code: code,
		Message: func() string {
			result := NewInfo(err)
			if codeMsg != "" && result != "" {
				result += ", " + codeMsg
			} else if codeMsg != "" {
				result = codeMsg
			}
			return strings.TrimSpace(result)
		}(),
	}
}
func NewInfo(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// SetRes Response res
func (g *Gin) SetRes(res interface{}, err error, code int) {
	g.JSON(http.StatusOK, NewRes(res, err, code))
}

// SetJson Set Json
func (g *Gin) SetJson(res interface{}) {
	g.SetRes(res, nil, codeSuccess)
}

// SetError Check Error
func (g *Gin) SetError(code int, err error) {
	g.SetRes(nil, err, code)
	g.Abort()
}

const (
	codeParamErr = iota + 1000
	codeRunning
	codePending
	codeExecErr
	codeNoData
	codeSuccess = 0
	codeErrApp  = 500
)

var MsgFlags = map[int]string{
	codeParamErr: "parameter error",
	codeRunning:  "is executing",
	codePending:  "waiting to run",
	codeExecErr:  "failed to execute",
	codeNoData:   "no data",
	codeSuccess:  "success",
	codeErrApp:   "internal error",
}

// GetMsg get error information based on Code
func getMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[codeErrApp]
}
