package utils

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const ServiceName = "osreapi"

var (
	LogLevel      string
	ListenAddress string
	MaxRequests   int64
	DataPath      string
	ExecTimeOut   time.Duration
	TmpDir        = filepath.Join(os.TempDir(), ServiceName)
)

func FileOrPathExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func EnsureDirExist(name string) error {
	if !FileOrPathExist(name) {
		return os.MkdirAll(name, os.ModePerm)
	}
	return nil
}

func GbkToUtf8(s []byte) []byte {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		logrus.Error(err)
		return s
	}
	return b
}

func ClearTmpDirOldScript() {
	ds, err := os.ReadDir(TmpDir)
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, d := range ds {
		if strings.HasSuffix(DataPath, d.Name()) && d.IsDir() {
			continue
		}
		if strings.HasSuffix(d.Name(), ".log") && !d.IsDir() {
			continue
		}
		_ = os.RemoveAll(filepath.Join(TmpDir, d.Name()))
	}
}
