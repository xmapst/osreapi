package pool

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/router/base"
	"github.com/xmapst/osreapi/internal/router/types"
	"github.com/xmapst/osreapi/internal/worker"
)

// Detail
// @Summary pool detail
// @description detail pool
// @Tags Pool
// @Accept application/json
// @Accept application/toml
// @Accept application/x-yaml
// @Accept multipart/form-data
// @Produce application/json
// @Produce application/x-yaml
// @Produce application/toml
// @Success 200 {object} types.BaseRes
// @Failure 500 {object} types.BaseRes
// @Router /api/v1/pool [get]
func Detail(c *gin.Context) {
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
	if ws == nil {
		render.SetRes(&types.Pool{
			Size:    worker.GetSize(),
			Total:   worker.GetTotal(),
			Running: worker.Running(),
			Waiting: worker.Waiting(),
		})
		return
	}
	// websocket 方式
	defer func() {
		_ = ws.WriteControl(websocket.CloseMessage, nil, time.Now())
		_ = ws.Close()
	}()
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		res := &types.Pool{
			Size:    worker.GetSize(),
			Total:   worker.GetTotal(),
			Running: worker.Running(),
			Waiting: worker.Waiting(),
		}
		err := ws.WriteJSON(base.NewRes(res, nil, base.CodeSuccess))
		if err != nil {
			logx.Errorln(err)
			return
		}
	}
}
