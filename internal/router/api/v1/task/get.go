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
	"github.com/xmapst/osreapi/internal/utils"
	"github.com/xmapst/osreapi/internal/worker"
)

// Get
// @Summary task detail
// @description detail task
// @Tags Task
// @Accept application/json
// @Accept application/toml
// @Accept application/x-yaml
// @Accept multipart/form-data
// @Produce application/json
// @Produce application/x-yaml
// @Produce application/toml
// @Param task path string true "task id"
// @Success 200 {object} types.BaseRes
// @Failure 500 {object} types.BaseRes
// @Router /api/v1/task/{task} [get]
func Get(c *gin.Context) {
	render := base.Gin{Context: c}
	var ws *websocket.Conn
	if websocket.IsWebSocketUpgrade(c.Request) {
		var err error
		ws, err = base.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logx.Errorln(err)
			render.SetError(base.CodeErrNoData, err)
			return
		}
	}

	if ws == nil {
		res, err, code := get(c)
		switch code {
		case base.CodePending, base.CodePaused:
			render.SetNegotiate(res, err, base.CodeRunning)
		default:
			render.SetNegotiate(res, err, code)
		}
		return
	}
	// websocket 方式
	defer func() {
		_ = ws.WriteControl(websocket.CloseMessage, nil, time.Now())
		_ = ws.Close()
	}()
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		res, err, code := get(c)
		if code == base.CodeErrNoData {
			err = ws.WriteJSON(base.NewRes(res, err, code))
			if err != nil {
				logx.Errorln(err)
			}
			return
		}
		err = ws.WriteJSON(base.NewRes(res, err, code))
		if err != nil {
			logx.Errorln(err)
			return
		}
	}
}

func get(c *gin.Context) ([]types.TaskDetailRes, error, int) {
	task := c.Param("task")
	if task == "" {
		return nil, errors.New("task does not exist"), base.CodeErrNoData
	}
	taskState, err := storage.TaskDetail(task)
	if err != nil {
		return nil, err, base.CodeErrNoData
	}
	state := exec.StateMap[taskState.State]
	c.Request.Header.Set(types.XTaskState, state)
	c.Writer.Header().Set(types.XTaskState, state)
	c.Set(types.XTaskState, state)
	tasksStepStates, err := storage.TaskStepList(task)
	if err != nil {
		return nil, err, base.CodeErrNoData
	}
	var scheme = "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	var stopMsg, runMsg, pausedMsg, pendingMsg []string
	var exitCode int64
	var res []types.TaskDetailRes
	var uriPrefix = fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, strings.TrimSuffix(c.Request.URL.Path, "/"))
	for _, v := range tasksStepStates {
		exitCode += v.Code
		msg := v.Message
		if v.Code != 0 {
			msg = fmt.Sprintf("Step: %s, Exit Code: %d", v.ID, v.Code)
			if taskState.MetaData.VMInstanceID != "" {
				msg += fmt.Sprintf(", Instance ID: %s", taskState.MetaData.VMInstanceID)
			}
			if taskState.MetaData.HardWareID != "" {
				msg += fmt.Sprintf(", Hardware ID: %s", taskState.MetaData.HardWareID)
			}
			stopMsg = append(stopMsg, msg)
		}
		var output []string
		if v.State == exec.Stop {
			outputs, _ := storage.TaskStepLogList(task, v.ID)
			for _, o := range outputs {
				if o.Content == worker.ConsoleStart || o.Content == worker.ConsoleDone {
					continue
				}
				output = append(output, o.Content)
			}
		} else {
			msg = fmt.Sprintf("Step: %s", v.ID)
			switch v.State {
			case exec.Running:
				runMsg = append(runMsg, msg)
				output = []string{"The step is running"}
			case exec.Paused:
				pausedMsg = append(pausedMsg, msg)
				output = []string{"The step is paused"}
			case exec.Pending:
				pendingMsg = append(pendingMsg, msg)
			default:
				logx.Debugln("unknown state", v.State)
			}
		}

		if output == nil {
			output = []string{v.Message}
		}
		_res := types.TaskDetailRes{
			Name:           v.ID,
			State:          v.State,
			Code:           v.Code,
			Manager:        fmt.Sprintf("%s/step/%s", uriPrefix, v.ID),
			Workspace:      fmt.Sprintf("%s/workspace", uriPrefix),
			Message:        strings.Join(output, "\n"),
			EnvVars:        v.EnvVars,
			Timeout:        v.Timeout.String(),
			DependsOn:      v.DependsOn,
			CommandType:    v.CommandType,
			CommandContent: v.CommandContent,
			Times: &types.StrTimes{
				ST: utils.TimeStr(v.Times.ST),
				ET: utils.TimeStr(v.Times.ET),
			},
		}
		res = append(res, _res)
	}

	switch taskState.State {
	// 运行结束
	case exec.Stop:
		if exitCode != 0 {
			return res, fmt.Errorf("[%s]", strings.Join(stopMsg, "; ")), base.CodeExecErr
		}
		return res, nil, base.CodeSuccess
	// 运行中, 排队中
	case exec.Running:
		var err error
		if runMsg != nil && pausedMsg != nil {
			err = fmt.Errorf("currently executing: [%s]; paused:[%s]", strings.Join(runMsg, "; "), strings.Join(pausedMsg, "; "))
		} else if runMsg != nil {
			err = fmt.Errorf("currently executing: [%s]", strings.Join(runMsg, "; "))
		} else if pausedMsg != nil {
			err = fmt.Errorf("paused:[%s]", strings.Join(pausedMsg, "; "))
		}
		return res, err, base.CodeRunning
	case exec.Pending:
		return res, nil, base.CodePending
	case exec.Paused:
		return res, nil, base.CodePaused
	case exec.SystemErr:
		return res, fmt.Errorf(taskState.Message), base.CodeExecErr
	default:
		return nil, errors.New("task does not exist"), base.CodeErrNoData
	}
}
