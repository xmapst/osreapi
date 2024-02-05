package step

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/xmapst/osreapi/internal/dag"
	"github.com/xmapst/osreapi/internal/exec"
	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/router/base"
	"github.com/xmapst/osreapi/internal/storage"
)

// Manager
// @Summary task step manager
// @description manager task step
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
// @Param action query string false "management action" Enums(paused,kill,pause,resume) default(paused)
// @Param duration query string false "how long to pause; if empty, manual continuation is required" default(1m)
// @Success 200 {object} types.BaseRes
// @Failure 500 {object} types.BaseRes
// @Router /api/v1/task/{task}/step/{step} [put]
func Manager(c *gin.Context) {
	render := base.Gin{Context: c}
	taskID := c.Param("task")
	if taskID == "" {
		render.SetError(base.CodeErrNoData, errors.New("task does not exist"))
		return
	}
	stepID := c.Param("step")
	if stepID == "" {
		render.SetError(base.CodeErrNoData, errors.New("step does not exist"))
		return
	}
	action := c.DefaultQuery("action", "paused")
	duration := c.Query("duration")
	manager, err := dag.VertexManager(taskID, stepID)
	if err != nil {
		logx.Errorln(err)
		render.SetError(base.CodeErrNoData, err)
		return
	}
	step, err := storage.TaskStepDetail(taskID, stepID)
	if err != nil {
		logx.Errorln(err)
		render.SetError(base.CodeErrNoData, err)
		return
	}
	switch action {
	case "kill":
		err = manager.Kill()
		if err == nil {
			step.State = exec.Killed
			err = storage.SetTaskStep(taskID, stepID, step)
		}
	case "pause":
		if step.State == exec.Running {
			render.SetError(base.CodeExecErr, dag.ErrRunning)
			return
		}
		if !manager.Paused() {
			_ = manager.Pause(duration)
			step.State = exec.Paused
			err = storage.SetTaskStep(taskID, stepID, step)
		}
	case "resume":
		if manager.Paused() {
			manager.Resume()
			step.State = exec.Pending
			err = storage.SetTaskStep(taskID, stepID, step)
		}
	}
	if err != nil {
		logx.Errorln(err)
		render.SetError(base.CodeExecErr, err)
		return
	}
	render.SetRes(nil)
}
