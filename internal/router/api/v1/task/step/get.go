package step

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/xmapst/osreapi/internal/exec"
	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/router/base"
	"github.com/xmapst/osreapi/internal/storage"
	"github.com/xmapst/osreapi/internal/storage/types"
	"github.com/xmapst/osreapi/internal/worker"
)

// Get
// @Summary task step detail
// @description detail task step
// @Tags Task
// @Accept application/json
// @Accept application/toml
// @Accept application/x-yaml
// @Accept multipart/form-data
// @Produce application/json
// @Produce application/x-yaml
// @Produce application/toml
// @Param task path string true "task id"
// @Param step path string true "step id"
// @Success 200 {object} types.BaseRes
// @Failure 500 {object} types.BaseRes
// @Router /api/v1/task/{task}/step/{step} [get]
func Get(c *gin.Context) {
	var render = base.Gin{Context: c}
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
	var latest int64 = -1
	if ws != nil {
		// 使用websocket方式
		stepDetailWebsocket(c, ws, &latest)
		return
	}
	var task = c.Param("task")
	if task == "" {
		render.SetError(base.CodeErrNoData, errors.New("task does not exist"))
		return
	}
	var step = c.Param("step")
	if step == "" {
		render.SetError(base.CodeErrNoData, errors.New("step does not exist"))
		return
	}
	taskStepState, err := storage.TaskStepDetail(task, step)
	if err != nil {
		render.SetError(base.CodeErrNoData, err)
		return
	}
	var res, _ = getTaskStepLog(task, step, &latest)
	switch taskStepState.State {
	case exec.Pending:
		render.SetNegotiate([]types.TaskStepLog{
			{
				Timestamp: time.Now().UnixNano(),
				Line:      1,
				Content:   taskStepState.Message,
			},
		}, nil, base.CodePending)
		return
	case exec.Paused:
		render.SetNegotiate([]types.TaskStepLog{
			{
				Timestamp: time.Now().UnixNano(),
				Line:      1,
				Content:   "step is paused",
			},
		}, nil, base.CodePaused)
		return
	case exec.Running:
		render.SetNegotiate(res, errors.New("in progress"), base.CodeRunning)
		return
	default:
		logx.Debugln("unknown state", taskStepState.State)
	}
	if taskStepState.Code != 0 {
		if res == nil {
			res = []*types.TaskStepLog{
				{
					Timestamp: time.Now().UnixNano(),
					Line:      0,
					Content:   taskStepState.Message,
				},
			}
		}
		err = fmt.Errorf("exit code: %d", taskStepState.Code)
		if taskStepState.Message != "" {
			err = fmt.Errorf(taskStepState.Message)
		}
		render.SetNegotiate(res, err, base.CodeExecErr)
		return
	}
	render.SetRes(res)
}

func stepDetailWebsocket(c *gin.Context, ws *websocket.Conn, latest *int64) {
	defer func() {
		_ = ws.WriteControl(websocket.CloseMessage, nil, time.Now())
		_ = ws.Close()
	}()
	var task = c.Param("task")
	if task == "" {
		err := ws.WriteJSON(base.NewRes(nil, errors.New("task does not exist"), base.CodeErrNoData))
		if err != nil {
			logx.Errorln(err)
		}
		return
	}
	var step = c.Param("step")
	if step == "" {
		err := ws.WriteJSON(base.NewRes(nil, errors.New("step does not exist"), base.CodeErrNoData))
		if err != nil {
			logx.Errorln(err)
		}
		return
	}

	var ticker = time.NewTicker(1 * time.Second)
	var pendingOnce sync.Once
	var pausedOnce sync.Once
	for range ticker.C {
		taskStepState, err := storage.TaskStepDetail(task, step)
		if err != nil {
			err = ws.WriteJSON(base.NewRes(nil, err, base.CodeErrNoData))
			if err != nil {
				logx.Errorln(err)
			}
			return
		}
		if taskStepState.State == exec.Running {
			break
		}
		switch taskStepState.State {
		case exec.Pending:
			// 只发送一次
			pendingOnce.Do(func() {
				err = ws.WriteJSON(base.NewRes([]types.TaskStepLog{
					{
						Timestamp: time.Now().UnixNano(),
						Line:      1,
						Content:   taskStepState.Message,
					},
				}, nil, base.CodePending))
				if err != nil {
					logx.Errorln(err)
					return
				}
			})
			continue
		case exec.Paused:
			// 只发送一次
			pausedOnce.Do(func() {
				err = ws.WriteJSON(base.NewRes([]types.TaskStepLog{
					{
						Timestamp: time.Now().UnixNano(),
						Line:      1,
						Content:   "step is paused",
					},
				}, nil, base.CodePaused))
				if err != nil {
					logx.Errorln(err)
					return
				}
			})
			continue
		default:
			res, _ := getTaskStepLog(task, step, latest)
			if res == nil {
				res = []*types.TaskStepLog{
					{
						Timestamp: time.Now().UnixNano(),
						Line:      0,
						Content:   taskStepState.Message,
					},
				}
			}
			if taskStepState.Code != 0 {
				errMsg := fmt.Errorf("exit code: %d", taskStepState.Code)
				if taskStepState.Message != "" {
					errMsg = fmt.Errorf(taskStepState.Message)
				}
				err = ws.WriteJSON(base.NewRes(res, errMsg, base.CodeExecErr))
				if err != nil {
					logx.Errorln(err)
				}
				return
			}
			err = ws.WriteJSON(base.NewRes(res, nil, base.CodeSuccess))
			if err != nil {
				logx.Errorln(err)
			}
			return
		}
	}

	for range ticker.C {
		res, done := getTaskStepLog(task, step, latest)
		for _, v := range res {
			if err := ws.WriteJSON(base.NewRes(v, errors.New("in progress"), base.CodeRunning)); err != nil {
				logx.Errorln(err)
				return
			}
		}

		if *latest <= 0 || done {
			return
		}
	}
}

func getTaskStepLog(task, step string, latest *int64) (res []*types.TaskStepLog, done bool) {
	outputs, err := storage.TaskStepLogList(task, step)
	if err != nil {
		return
	}
	for k, v := range outputs {
		if v.Content == worker.ConsoleStart {
			continue
		}
		if v.Content == worker.ConsoleDone {
			done = true
			continue
		}
		if int64(k) <= *latest {
			continue
		}
		res = append(res, v)
	}
	*latest = int64(len(outputs)) - 1
	return
}
