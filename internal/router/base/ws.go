package base

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/xmapst/osreapi/internal/logx"
)

var upGrader = websocket.Upgrader{
	HandshakeTimeout: 3 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MsgHandler func(messageType int, p []byte) error

func Upgrade(w http.ResponseWriter, r *http.Request, handler MsgHandler) (*websocket.Conn, error) {
	ws, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	// other operations
	// ws.SetPongHandler(func(appData string) error { return nil })
	// ws.SetPingHandler(func(appData string) error { return nil })
	// ...
	ws.SetCloseHandler(func(code int, text string) error {
		return nil
	})
	go readMessage(ws, handler)
	return ws, nil
}

func readMessage(ws *websocket.Conn, handler MsgHandler) {
	defer func() {
		if ws == nil {
			return
		}
		_ = ws.WriteControl(websocket.CloseMessage, nil, time.Now())
		_ = ws.Close()
	}()
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			logx.Errorln(err)
			return
		}
		switch messageType {
		case websocket.CloseMessage:
			logx.Infoln("close websocket connect", ws.RemoteAddr().String())
			return
		case websocket.PingMessage:
			err = ws.WriteMessage(websocket.PongMessage, []byte("ping"))
			if err != nil {
				logx.Warningf("send %s ping err: %v", ws.RemoteAddr().String(), err)
			}
		case websocket.PongMessage:
			err = ws.WriteMessage(websocket.PongMessage, []byte("pong"))
			if err != nil {
				logx.Warningf("send %s pong err: %v", ws.RemoteAddr().String(), err)
			}
		default:
			if handler != nil {
				err = handler(messageType, p)
				if err != nil {
					logx.Errorln(err)
				}
			}
		}
	}
}
