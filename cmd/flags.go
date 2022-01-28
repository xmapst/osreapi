package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/xmapst/osreapi/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"path/filepath"
)

func init() {
	// flags
	kingpin.Flag(
		"addr",
		`host:port for execution.`,
	).Default(":2376").StringVar(&utils.ListenAddress)
	kingpin.Flag(
		"log_level",
		"Set log output level.",
	).Default("info").StringVar(&utils.LogLevel)
	kingpin.Flag(
		"exec_timeout",
		`Set the default exec command expire time.`,
	).Default("30m").DurationVar(&utils.ExecTimeOut)
	kingpin.Flag(
		"max-requests",
		`Maximum number of concurrent requests. 0 to disable.`,
	).Default("0").Int64Var(&utils.MaxRequests)
	kingpin.Flag(
		"data",
		"Save data with memory or custom path.",
	).Default(filepath.Join(utils.TmpDir, "data")).StringVar(&utils.DataPath)
	_ = utils.EnsureDirExist(utils.TmpDir)
	// log format init
	logrus.SetReportCaller(true)
	logrus.AddHook(utils.NewRotateHook())
	logrus.SetFormatter(&utils.FileFormatter{})
}
