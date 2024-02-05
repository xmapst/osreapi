package server

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/osext"
	"github.com/xmapst/osreapi/internal/selfupdate"
	"github.com/xmapst/osreapi/internal/worker"
)

func (p *Program) selfUpdate() {
	_, err := p.cron.AddFunc("@hourly", func() {
		p.doUpdate()
	})
	if err != nil {
		logx.Warnln(err)
		return
	}
}

func (p *Program) doUpdate() {
	if p.sURL == "" {
		logx.Debugln("automatic updates are not turned on")
		return
	}

	if p.sHash == nil {
		p.localSha256sum()
	}

	// 获取远端hash, example: https://down.q1.com/autoops/tools/windows_remote_executor/latest.sha256sum
	checksum, name := p.remoteSha256sum(fmt.Sprintf("%s/latest.sha256sum", p.sURL))
	if checksum == nil || p.sHash == nil || bytes.Equal(p.sHash, checksum) {
		// 获取不到本地, 远端hash, 或者两者相同不更新
		return
	}

	opts := selfupdate.Options{
		Checksum: checksum,
	}
	// 检查权限
	if err := opts.CheckPermissions(); err != nil {
		logx.Warnln(err)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/%s", p.sURL, name))
	if resp != nil {
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				logx.Errorln(err)
			}
		}(resp.Body)
	}
	if err != nil {
		logx.Warnln(err)
		return
	}

	err = selfupdate.PrepareAndCheckBinary(resp.Body, opts)
	if err != nil {
		logx.Warnln(err)
		return
	}

	if (worker.Running() + worker.Waiting()) != 0 {
		// 还有任务执行中或者等待执行不升级
		logx.Warnln("the task has not been completed")
		return
	}

	// 关闭入口
	p.close()

	// 更新文件
	err = selfupdate.CommitBinary(opts)
	if err != nil {
		logx.Errorln(err)
		// 重新启动http以及tcp入口
		err = p.startServer()
		if err == nil {
			return
		}
	}

	// 依赖systemd或者service守护
	os.Exit(200)
}

func (p *Program) localSha256sum() {
	path, err := osext.Executable()
	if err != nil {
		logx.Warnln(err)
		return
	}
	hash := sha256.New()
	file, err := os.Open(path)
	if err != nil {
		logx.Warnln(err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logx.Errorln(err)
		}
	}(file)
	_, err = io.Copy(hash, file)
	if err != nil {
		logx.Warnln(err)
		return
	}
	p.sHash = hash.Sum([]byte{})
	return
}

func (p *Program) remoteSha256sum(url string) ([]byte, string) {
	resp, err := http.Get(url)
	if resp != nil {
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				logx.Errorln(err)
			}
		}(resp.Body)
	}
	if err != nil {
		logx.Warnln(err)
		return nil, ""
	}
	body := bufio.NewReader(resp.Body)
	for {
		line, _, err := body.ReadLine()
		if err == io.EOF {
			return nil, ""
		}
		if err != nil {
			logx.Warnln(err)
		}
		if bytes.Contains(line, []byte(runtime.GOOS)) && bytes.Contains(line, []byte(runtime.GOARCH)) {
			fields := bytes.Fields(line)
			if len(fields) != 2 {
				continue
			}
			src := fields[0]
			n, err := hex.Decode(src, src)
			if err == nil {
				return src[:n], filepath.Base(string(fields[1]))
			} else {
				logx.Warnln(err)
			}
		}
	}
}
