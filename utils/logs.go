package utils

import (
	"fmt"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/utahta/go-cronowriter"
	"os"
	"path"
	"strings"
)

const logName = ServiceName + ".log"

func NewRotateHook() *lfshook.LfsHook {
	if len(TmpDir) != 0 {
		err := os.MkdirAll(TmpDir, 0777)
		if err != nil {
			logrus.Errorf("config local file system logger error. %+v", err)
		}
	}
	baseLogPath := path.Join(TmpDir, logName)
	pattern := path.Join(TmpDir, "%Y-%m-%d_"+logName)
	writer := cronowriter.MustNew(pattern,
		cronowriter.WithSymlink(baseLogPath),
	)
	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &FileFormatter{})
}

type FileFormatter struct {
	logrus.TextFormatter
}

func (c *FileFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	logStr := fmt.Sprintf("%s %s %s:%d %v",
		entry.Time.Format("2006/01/02 15:04:05"),
		strings.ToUpper(entry.Level.String()),
		path.Base(entry.Caller.File),
		entry.Caller.Line,
		entry.Message,
	)
	if len(entry.Data) != 0 {
		logStr += fmt.Sprintf("; %v", entry.Data)
	}
	return []byte(logStr + "\n"), nil
}
