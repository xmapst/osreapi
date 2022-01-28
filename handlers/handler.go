package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/osreapi/cache"
	"github.com/xmapst/osreapi/engine"
	"github.com/xmapst/osreapi/utils"
	"net/http"
	"strings"
	"time"
)

type listRes struct {
	ID    string `json:"id,omitempty"`
	State string `json:"state,omitempty"`
	Times *times `json:"times,omitempty"`
}

type times struct {
	Start     string `json:"start,omitempty"`
	Completed string `json:"completed,omitempty"`
	Expired   string `json:"expired,omitempty"`
}

// List
// @Summary List
// @description List all current task IDs
// @Tags Exec
// @Param sort query string false "sort by [start,ttl,completed]"
// @Success 200 {object} JSONResult{}
// @Failure 500 {object} JSONResult{}
// @Router / [get]
func List(c *gin.Context) {
	render := Gin{Context: c}
	var _res []*cache.ListData
	sortBy := c.DefaultQuery("sort", "start")
	switch strings.ToLower(sortBy) {
	case "start":
		_res = cache.GetAllByStartTimes()
	case "ttl":
		_res = cache.GetAllByTTL()
	case "completed":
		_res = cache.GetAllByEndTimes()
	default:
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":    codeParamErr,
			"message": "sort parameter error",
		})
		return
	}
	if len(_res) == 0 {
		render.SetJson(nil)
		return
	}
	var res []*listRes
	for _, v := range _res {
		_res := &listRes{
			ID:    v.ID,
			State: cache.StateMap[v.State],
			Times: new(times),
		}
		if v.StartTimes != 0 {
			_res.Times.Start = time.Unix(0, v.StartTimes).Format(time.RFC3339)
		}
		if v.CompletedTimes != 0 {
			_res.Times.Completed = time.Unix(0, v.CompletedTimes).Format(time.RFC3339)
		}
		if v.TTL != 0 {
			_res.Times.Expired = v.TTL.String()
		}
		res = append(res, _res)
	}
	render.SetJson(map[string]interface{}{
		"total":   len(res),
		"tasks":   res,
		"running": engine.Pool.QueueLength(),
	})
}

// Get
// @Summary Get
// @description Get the execution result
// @Tags Exec
// @Param id path string true "id"
// @Success 200 {object} JSONResult{}
// @Failure 500 {object} JSONResult{}
// @Router /{id} [get]
func Get(c *gin.Context) {
	render := Gin{Context: c}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":    codeParamErr,
			"message": "missing id parameter",
		})
		return
	}
	found, val := cache.Get(id)
	if !found {
		render.SetError(codeNoData, errors.New("id does not exist"))
		return
	}
	if val.State == cache.Stopped {
		for _, v := range val.Data {
			// exec error
			if v.ExitCode != 0 {
				render.SetRes(val.Data, fmt.Errorf("%s %d", v.Name, v.Step), codeExecErr)
				return
			}
		}
		render.SetJson(val.Data)
		return
	}
	if val.State == cache.Pending {
		render.SetError(codePending, nil)
		return
	}
	render.SetError(codeRunning, nil)
}

const xTaskID = "X-Task-ID"

type PostStruct struct {
	Name           string            `json:"name" form:"name" description:"名称[可选], 不可重复"`
	CommandType    string            `json:"command_type" form:"command_type" binding:"required" description:"命令类型[必选]" example:"sh"`
	CommandContent string            `json:"command_content" form:"command_content" binding:"required" description:"命令内容[必选]" example:"ping baidu.com -c 10"`
	DependsOn      []string          `json:"depends_on,omitempty" form:"depends_on,omitempty" description:"依赖步骤[可选]" example:""`
	Envs           map[string]string `json:"envs,omitempty" form:"envs,omitempty" description:"外部环境变量[可选]" example:"env1:value1,env2:value2"`
	Timeout        string            `json:"timeout,omitempty" form:"timeout,omitempty" description:"执行超时时间[可选]" example:"5m"`
	timeout        time.Duration
}

type PostRes struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	StepCount int    `json:"step_count"`
}

// Post
// @Summary Exec
// @description Execute command or script content
// @Tags Exec
// @Param scripts body []PostStruct true "scripts"
// @Success 200 {object} JSONResult{}
// @Failure 500 {object} JSONResult{}
// @Router / [post]
func Post(c *gin.Context) {
	render := Gin{Context: c}
	var req []*PostStruct
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	if len(req) == 0 {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":    codeParamErr,
			"message": getMsg(codeParamErr),
		})
		return
	}

	// check timeout
	var totalTimeOut time.Duration
	var err error
	for k, v := range req {
		var _timeout time.Duration
		if v.Timeout == "" {
			_timeout = utils.ExecTimeOut
		} else {
			_timeout, err = time.ParseDuration(v.Timeout)
			if err != nil {
				c.JSON(http.StatusOK, map[string]interface{}{
					"code":    http.StatusBadRequest,
					"message": fmt.Sprintf("step %d %v", k, err),
				})
				return
			}
		}
		req[k].timeout = _timeout
		totalTimeOut += _timeout
	}

	// create exec id
	taskID := uuid.NewV4().String()
	logrus.Info("new task ", taskID)

	// set header
	c.Request.Header.Set(xTaskID, taskID)
	c.Writer.Header().Set(xTaskID, taskID)
	c.Set(xTaskID, taskID)
	var jobs []*cache.Jobs
	for k, v := range req {
		if v.Name == "" {
			v.Name = fmt.Sprintf("%s_%d", taskID, k)
		} else {
			// Prevent name collisions
			v.Name = fmt.Sprintf("%s_%d", v.Name, k)
		}
		jobs = append(jobs, &cache.Jobs{
			Name:           v.Name,
			CommandType:    v.CommandType,
			CommandContent: v.CommandContent,
			EnvVars:        utils.Slice(v.Envs),
			ExecTimeout:    v.timeout,
			DependsOn:      v.DependsOn,
		})
	}
	// Join the pool for asynchronous processing
	engine.Process(taskID, jobs, totalTimeOut)
	render.SetJson(PostRes{
		ID:        taskID,
		Timestamp: time.Now().Unix(),
		StepCount: len(req),
	})
}
