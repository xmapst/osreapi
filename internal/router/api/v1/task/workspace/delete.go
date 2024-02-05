package workspace

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/router/base"
	"github.com/xmapst/osreapi/internal/server/config"
	"github.com/xmapst/osreapi/internal/utils"
)

// Delete
// @Summary task workspace delete
// @description task workspace delete
// @Tags Task
// @Accept application/json
// @Accept application/toml
// @Accept application/x-yaml
// @Accept multipart/form-data
// @Produce application/json
// @Produce application/x-yaml
// @Produce application/toml
// @Param task path string true "task id"
// @Param path query string false "dir path"
// @Success 200 {object} types.BaseRes
// @Failure 500 {object} types.BaseRes
// @Router /api/v1/task/{task}/workspace [delete]
func Delete(c *gin.Context) {
	render := base.Gin{Context: c}
	task := c.Param("task")
	if task == "" {
		render.SetError(base.CodeErrNoData, errors.New("task does not exist"))
		return
	}
	prefix := filepath.Join(config.App.WorkSpace, task)
	if !utils.FileOrPathExist(prefix) {
		render.SetError(base.CodeErrNoData, errors.New("task does not exist"))
		return
	}
	path := filepath.Join(prefix, utils.PathEscape(c.Query("path")))
	if err := os.RemoveAll(path); err != nil {
		logx.Errorln(err)
		render.SetError(base.CodeExecErr, err)
		return
	}
	render.SetRes(nil)
}
