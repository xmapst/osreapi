package base

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/xmapst/osreapi/internal/router/types"
)

const (
	CodeSuccess = 0
	CodeRunning = iota + 1000
	CodeExecErr
	CodeErrNoData
	CodePending
	CodePaused
)

var MsgFlags = map[int]string{
	CodeRunning:   "running",
	CodeExecErr:   "exec error",
	CodeErrNoData: "no data",
	CodePending:   "pending",
	CodePaused:    "paused",
	CodeSuccess:   "success",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[CodeErrNoData]
}

type Gin struct {
	*gin.Context
}

func NewRes(data interface{}, err error, code int) *types.BaseRes {
	if code == 200 {
		code = 0
	}
	return &types.BaseRes{
		Data:      data,
		Code:      code,
		Timestamp: time.Now().UnixNano(),
		Message: func() string {
			result := NewInfo(err)
			if result == "" {
				result = GetMsg(code)
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

// SetNegotiate Response res
func (g *Gin) SetNegotiate(res interface{}, err error, code int) {
	switch g.NegotiateFormat(binding.MIMEJSON, binding.MIMEYAML, binding.MIMETOML) {
	case binding.MIMEJSON:
		g.JSON(http.StatusOK, NewRes(res, err, code))

	case binding.MIMEYAML:
		g.YAML(http.StatusOK, NewRes(res, err, code))

	case binding.MIMETOML:
		g.TOML(http.StatusOK, NewRes(res, err, code))

	default:
		g.JSON(http.StatusOK, NewRes(res, err, code))
	}
}

// SetRes res
func (g *Gin) SetRes(res interface{}) {
	g.SetNegotiate(res, nil, CodeSuccess)
}

// SetError Check Error
func (g *Gin) SetError(code int, err error) {
	g.SetNegotiate(nil, err, code)
	g.Abort()
}
