package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/router/base"
	"github.com/xmapst/osreapi/internal/router/types"
)

func healthyz(c *gin.Context) {
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
		render.SetRes(&types.Healthyz{
			Server: c.Request.Host,
			Client: c.ClientIP(),
			State:  "Running",
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
		res := &types.Healthyz{
			Server: c.Request.Host,
			Client: c.ClientIP(),
			State:  "Running",
		}
		err := ws.WriteJSON(base.NewRes(res, nil, base.CodeSuccess))
		if err != nil {
			logx.Errorln(err)
			return
		}
	}
}
