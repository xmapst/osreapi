package zap

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xmapst/osreapi/internal/logx"
)

func Logger(c *gin.Context) {
	start := time.Now().UTC()
	path := c.Request.URL.String()
	c.Next()
	end := time.Now().UTC()
	latency := end.Sub(start)

	if len(c.Errors) > 0 {
		for _, e := range c.Errors.Errors() {
			logx.Errorln(c.ClientIP(), c.Request.Method, c.Request.Proto, c.Writer.Status(), path, latency, e)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		logx.Infoln(c.ClientIP(), c.Request.Method, c.Request.Proto, c.Writer.Status(), path, latency, c.Request.UserAgent())
	}
}

func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// Check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			httpRequest, _ := httputil.DumpRequest(c.Request, false)
			if brokenPipe {
				logx.Errorln(c.Request.URL.Path, httpRequest, err)
				// If the connection is dead, we can't write a status to it.
				_ = c.Error(err.(error)) // nolint: errcheck
				c.Abort()
				return
			}

			logx.Errorln("[Recovery from panic]", time.Now().UTC().Format(time.RFC3339), string(httpRequest), string(debug.Stack()), err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	c.Next()
}
