package config

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/xmapst/osreapi/internal/logx"
)

var App = new(Config)

type Config struct {
	ServiceName   string
	Debug         bool
	Normal        bool
	ListenAddress string
	MaxRequests   int64
	PoolSize      int
	ExecTimeOut   time.Duration
	WebTimeout    time.Duration
	ScriptDir     string
	LogDir        string
	WorkSpace     string
	DataDir       string
	DBType        string
	SelfUpdateURL string
}

func init() {
	executable, err := os.Executable()
	if err != nil {
		logx.Fatalln(err)
	}
	App.ServiceName = strings.TrimSuffix(filepath.Base(executable), ".exe")
}

func (c *Config) Init() error {
	c.LogDir = c.dir(c.LogDir, "logs")
	c.ScriptDir = c.dir(c.ScriptDir, "scripts")
	c.WorkSpace = c.dir(c.WorkSpace, "workspace")
	c.DataDir = c.dir(c.DataDir, "data")
	if c.DBType == "" {
		c.DBType = "bolt"
	}
	return nil
}

func (c *Config) dir(dir, sub string) string {
	dir = os.Expand(dir, func(s string) string {
		if s == "TMP" || s == "TEMP" || s == "TEMPDIR" || s == "TMPDIR" {
			return os.TempDir()
		}
		if s == "HOME" || s == "HOMEDIR" {
			_s, _ := os.UserHomeDir()
			return _s
		}
		return os.Getenv(s)
	})
	defer func() {
		if err := os.MkdirAll(dir, os.ModeDir); err != nil {
			logx.Fatalln(err)
		}
		if sub == "logs" {
			var logfile string
			if !c.Debug {
				logfile = filepath.Join(c.LogDir, c.ServiceName+".log")
			}
			logx.SetupLogger(logfile, zap.AddStacktrace(zapcore.ErrorLevel))
		}
		logx.Infoln(sub, "dir", dir)
	}()
	if dir != "" {
		return dir
	}
	dir = filepath.Join(os.TempDir(), c.ServiceName, sub)
	return dir
}
