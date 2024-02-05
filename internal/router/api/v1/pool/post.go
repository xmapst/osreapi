package pool

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/router/base"
	"github.com/xmapst/osreapi/internal/router/types"
	"github.com/xmapst/osreapi/internal/worker"
)

// Post
// @Summary pool setting
// @description post task step
// @Tags Pool
// @Accept application/json
// @Accept application/toml
// @Accept application/x-yaml
// @Accept multipart/form-data
// @Produce application/json
// @Produce application/x-yaml
// @Produce application/toml
// @Param setting body types.Pool true "pool setting"
// @Success 200 {object} types.BaseRes
// @Failure 500 {object} types.BaseRes
// @Router /api/v1/pool [post]
func Post(c *gin.Context) {
	render := base.Gin{Context: c}
	var req types.Pool
	if err := c.ShouldBind(&req); err != nil {
		logx.Errorln(err)
		render.SetError(base.CodeExecErr, err)
		return
	}
	if (worker.Running() != 0 || worker.Waiting() != 0) && req.Size <= worker.GetSize() {
		render.SetError(base.CodeExecErr, errors.New("there are still tasks running, scaling down is not allowed"))
		return
	}
	worker.SetSize(req.Size)
	render.SetRes(&types.Pool{
		Size:    worker.GetSize(),
		Running: worker.Running(),
		Waiting: worker.Waiting(),
	})
}
