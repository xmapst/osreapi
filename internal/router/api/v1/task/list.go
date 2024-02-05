package task

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/xmapst/osreapi/internal/exec"
	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/router/base"
	"github.com/xmapst/osreapi/internal/router/types"
	"github.com/xmapst/osreapi/internal/storage"
	"github.com/xmapst/osreapi/internal/storage/backend"
	"github.com/xmapst/osreapi/internal/utils"
)

// List
// @Summary task list
// @description get all task
// @Tags Task
// @Accept application/json
// @Accept application/toml
// @Accept application/x-yaml
// @Accept multipart/form-data
// @Produce application/json
// @Produce application/x-yaml
// @Produce application/toml
// @Success 200 {object} types.BaseRes
// @Failure 500 {object} types.BaseRes
// @Router /api/v1/task [get]
func List(c *gin.Context) {
	render := base.Gin{Context: c}
	var ws *websocket.Conn
	if websocket.IsWebSocketUpgrade(c.Request) {
		var err error
		ws, err = base.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logx.Errorln(err)
			render.SetNegotiate(&types.TaskListRes{}, err, base.CodeSuccess)
			return
		}
	}
	if ws == nil {
		res, err := list(c)
		if err != nil && !errors.Is(err, backend.ErrNotExist) {
			logx.Errorln(err)
		}
		render.SetRes(res)
		return
	}
	// websocket 方式
	defer func() {
		_ = ws.WriteControl(websocket.CloseMessage, nil, time.Now())
		_ = ws.Close()
	}()
	// 使用websocket方式
	var ticker = time.NewTicker(1 * time.Second)
	for range ticker.C {
		res, err := list(c)
		if err != nil {
			if !errors.Is(err, backend.ErrNotExist) {
				logx.Errorln(err)
			}
			err = ws.WriteJSON(base.NewRes(res, err, base.CodeSuccess))
			if err != nil {
				logx.Errorln(err)
			}
			return
		}
		err = ws.WriteJSON(base.NewRes(res, nil, base.CodeSuccess))
		if err != nil {
			logx.Errorln(err)
			return
		}
	}
}

// 每次获取全量数据
func list(c *gin.Context) (*types.TaskListRes, error) {
	tasksStates, err := storage.TaskList()
	if err != nil || len(tasksStates) == 0 {
		return nil, err
	}
	var res = &types.TaskListRes{
		Total: len(tasksStates),
	}
	var scheme = "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	var uriPrefix = fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, strings.TrimSuffix(c.Request.URL.Path, "/"))
	for _, v := range tasksStates {
		_res := &types.TaskList{
			ID:        v.ID,
			State:     v.State,
			Manager:   fmt.Sprintf("%s/%s", uriPrefix, v.ID),
			Workspace: fmt.Sprintf("%s/%s/workspace", uriPrefix, v.ID),
			EnvVars:   v.EnvVars,
			Timeout:   v.Timeout.String(),
			Count:     v.Count,
			Times: types.StrTimes{
				ST: utils.TimeStr(v.Times.ST),
				ET: utils.TimeStr(v.Times.ET),
			},
		}
		if v.State == exec.SystemErr {
			_res.Code = exec.SystemErr
			_res.Message = v.Message
		} else {
			tasksStepStates, err := storage.TaskStepList(v.ID)
			if err != nil {
				logx.Errorln(err)
				continue
			}
			var running, paused, pending, errorStateMsg []string
			var code int64
			for _, vv := range tasksStepStates {
				var state = fmt.Sprintf("Step: %s", vv.ID)
				if vv.State == exec.Running {
					running = append(running, state)
				}
				if vv.State == exec.Paused {
					paused = append(paused, state)
				}
				if vv.State == exec.Pending {
					pending = append(pending, state)
				}
				if vv.Code != 0 {
					errorStateMsg = append(errorStateMsg, state)
				}
				code += vv.Code
			}
			_res.Code = code
			switch {
			case v.State == exec.Stop || v.State == exec.Killed || v.State == exec.Timeout || v.State == exec.SystemErr:
				_res.Message = fmt.Sprintf("execution failed: [%s]", strings.Join(errorStateMsg, "; "))
				if _res.Code == 0 {
					_res.Message = "all steps executed successfully"
				}
			case v.State == exec.Running:
				if running != nil {
					_res.Message = fmt.Sprintf(
						"currently executing: [%s] ",
						strings.Join(running, "; "),
					)
				}
				if paused != nil {
					_res.Message += fmt.Sprintf(
						"paused:[%s]",
						strings.Join(paused, "; "),
					)
				}
				if pending != nil {
					_res.Message += fmt.Sprintf(
						"pending:[%s]",
						strings.Join(pending, "; "),
					)
				}
			case v.State == exec.Paused:
				_res.Message = "task is paused"
			default:
				if v.Message != "" {
					_res.Message = v.Message
				}
			}
		}
		res.Tasks = append(res.Tasks, _res)
	}
	return res, nil
}
