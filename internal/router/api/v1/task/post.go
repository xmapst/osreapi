package task

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/router/base"
	"github.com/xmapst/osreapi/internal/router/types"
	"github.com/xmapst/osreapi/internal/worker"
)

// Post
// @Summary post task
// @description post task step
// @Tags Task
// @Accept application/json
// @Accept application/toml
// @Accept application/x-yaml
// @Accept multipart/form-data
// @Produce application/json
// @Produce application/x-yaml
// @Produce application/toml
// @param id query string false "task id"
// @Param async query bool false "task asynchronously" default(false)
// @Param timeout query string false "task timeout"
// @Param env_vars query []string false "task envs"
// @Param scripts body types.Steps true "scripts"
// @Success 200 {object} types.BaseRes
// @Failure 500 {object} types.BaseRes
// @Router /api/v1/task [post]
func Post(c *gin.Context) {
	render := base.Gin{Context: c}
	var req = new(types.Task)
	if err := c.ShouldBindQuery(req); err != nil {
		logx.Errorln(err)
		render.SetError(base.CodeExecErr, err)
		return
	}

	if err := c.ShouldBind(&req.Step); err != nil {
		logx.Errorln(err)
		render.SetError(base.CodeExecErr, err)
		return
	}

	if err := req.Check(); err != nil {
		render.SetError(base.CodeExecErr, err)
		return
	}

	var task = &worker.Task{
		ID:       req.ID,
		Timeout:  req.TimeoutDuration,
		EnvVars:  req.EnvVars,
		MetaData: req.Step.GetMetaData(),
	}

	for _, v := range req.Step {
		task.Steps = append(task.Steps, &worker.TaskStep{
			ID:             v.Name,
			CommandType:    v.CommandType,
			CommandContent: v.CommandContent,
			EnvVars:        v.EnvVars,
			DependsOn:      v.DependsOn,
			Timeout:        v.TimeoutDuration,
		})
	}

	// 提交任务
	if err := worker.Submit(task); err != nil {
		logx.Errorln(err)
		render.SetError(base.CodeExecErr, err)
		return
	}

	c.Request.Header.Set(types.XTaskID, task.ID)
	c.Writer.Header().Set(types.XTaskID, task.ID)
	c.Set(types.XTaskID, task.ID)

	var scheme = "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	render.SetRes(&types.TaskRes{
		URL:   fmt.Sprintf("%s://%s%s/%s", scheme, c.Request.Host, strings.TrimSuffix(c.Request.URL.Path, "/"), req.ID),
		ID:    req.ID,
		Count: len(req.Step),
	})
}
