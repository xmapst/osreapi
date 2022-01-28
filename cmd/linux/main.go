//go:build !windows

package main

import (
    "context"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "github.com/xmapst/osreapi"
    "github.com/xmapst/osreapi/cache"
    _ "github.com/xmapst/osreapi/cmd"
    "github.com/xmapst/osreapi/engine"
    "github.com/xmapst/osreapi/routers"
    "github.com/xmapst/osreapi/utils"
    "gopkg.in/alecthomas/kingpin.v2"
    "io/ioutil"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

// @title OSRemoteExecution API
// @version v1.0.0
// @description This is a os remote executor orchestration script interface.
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	kingpin.Version(osreapi.VersionIfo())
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
    level,err := logrus.ParseLevel(utils.LogLevel)
    if err != nil {
        logrus.Fatalln(err)
    }
    logrus.SetLevel(level)
    osreapi.PrintHeadInfo()
    if level != logrus.DebugLevel {
        gin.SetMode(gin.ReleaseMode)
        logrus.SetOutput(ioutil.Discard)
    }else {
        gin.SetMode(gin.DebugMode)
    }

    logrus.Info("starting, please wait")
    fmt.Println("starting, please wait")
	// clear old script
	utils.ClearTmpDirOldScript()
    engine.NewExecPool()
    cache.Init()

    // Resume a script that was not run last time
    logrus.Info("restoring a task that did not run last time")
    fmt.Println("restoring a task that did not run last time")
    engine.Reload()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	gin.DisableConsoleColor()
	srv := &http.Server{
		Addr:         utils.ListenAddress,
		WriteTimeout: 600 * time.Second,
		ReadTimeout:  600 * time.Second,
		IdleTimeout:  600 * time.Second,
		Handler:      routers.Router(level == logrus.DebugLevel),
	}
	go func() {
		logrus.Infof("listen address [%s]", utils.ListenAddress)
		if err := srv.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}()
    logrus.Info("server started")
    fmt.Println("server started")
	<-signals
	logrus.Info("shutdown server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_ = srv.Shutdown(ctx)
    cache.Close()
}
