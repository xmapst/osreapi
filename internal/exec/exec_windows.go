//go:build windows

package exec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"

	"github.com/go-cmd/cmd"
	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var acp = windows.GetACP()

func (c *Cmd) selfScriptSuffix() string {
	switch c.shell {
	case "cmd", "bat":
		return ".bat"
	case "powershell", "ps", "ps1":
		return ".ps1"
	default:
		return ".bat"
	}
}

func (c *Cmd) beforeExec() []func(cmd *exec.Cmd) {
	return []func(cmd *exec.Cmd){
		func(cmd *exec.Cmd) {
			cmd.SysProcAttr = &syscall.SysProcAttr{
				HideWindow: true,
			}
		},
	}
}

func (c *Cmd) selfCmd() {
	switch c.shell {
	case "cmd", "bat":
		c.cmd = cmd.NewCmdOptions(c.ops, "cmd", "/D", "/E:ON", "/V:OFF", "/S", "/C", c.scriptName)
	case "powershell", "ps", "ps1":
		// 解决用户不写exit时, powershell进程外获取不到退出码
		command := fmt.Sprintf("$ErrorActionPreference='Continue';%s;exit $LASTEXITCODE", c.scriptName)
		// 激进方式, 强制用户脚本没问题
		// command := fmt.Sprintf("$ErrorActionPreference='Stop';%s;exit $LASTEXITCODE", c.absFilePath)
		c.cmd = cmd.NewCmdOptions(c.ops, "powershell", "-NoLogo", "-NonInteractive", "-Command", command)
	default:
		c.cmd = cmd.NewCmdOptions(c.ops, "cmd", "/D", "/E:ON", "/V:OFF", "/S", "/C", c.scriptName)
	}
}

func (c *Cmd) Run(ctx context.Context) (code int64, err error) {
	defer func() {
		if _err := recover(); _err != nil {
			err = fmt.Errorf("%v", _err)
		}
	}()

	c.newCmd()
	defer func() {
		c.clear()
		c.cancel()
		if err != nil {
			c.logger.Errorln(c.scriptName, "exit code", code, err)
			return
		}
		c.logger.Infoln(c.scriptName, "exit code", code)
	}()

	// Print STDOUT and STDERR lines streaming from Cmd
	go c.consoleOutput()

	select {
	// 人工强制终止
	case <-ctx.Done():
		if c.cmd.Status().PID != 0 {
			if err = KillAll(c.cmd.Status().PID); err != nil {
				c.logger.Errorln(c.scriptName, err)
			}
		}
		err = ErrManual
		code = Killed
		if context.Cause(ctx) != nil {
			switch {
			case errors.Is(context.Cause(ctx), ErrTimeOut):
				err = ErrTimeOut
				code = Timeout
			default:
				err = context.Cause(ctx)
			}
		}
	// 执行超时信号
	case <-c.ctx.Done():
		// 如果直接使用cmd.Process.Kill()并不能杀死主进程下的所有子进程
		// _ = cmd.Process.Kill()
		if c.cmd.Status().PID != 0 {
			if err = KillAll(c.cmd.Status().PID); err != nil {
				c.logger.Errorln(c.scriptName, err)
			}
		}
		err = ErrTimeOut
		code = Timeout
	// 执行结果
	case status := <-c.cmd.Start():
		code = int64(status.Exit)
		err = status.Error
		if err != nil && code == 0 {
			code = SystemErr
		}
		if err == nil && code != 0 {
			err = fmt.Errorf("%s", last(c.stderrBuf.Lines()))
		}
		if err != nil {
			c.logger.Errorln(c.scriptName, err)
		}
	}
	return
}

func (c *Cmd) transform(line string) string {
	if c.isGBK(line) || acp == 936 {
		line = string(c.gbkToUtf8([]byte(line)))
	}
	return line
}

func (c *Cmd) gbkToUtf8(s []byte) []byte {
	defer func() {
		err := recover()
		if err != nil {
			c.logger.Errorln(c.workspace, c.scriptDir, c.scriptName, c.shell, err)
		}
	}()
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	b, err := io.ReadAll(reader)
	if err != nil {
		c.logger.Errorln(c.scriptName, err)
		return s
	}
	return b
}

func (c *Cmd) utf8ToGb2312(s string) string {
	defer func() {
		err := recover()
		if err != nil {
			c.logger.Errorln(c.scriptName, err)
		}
	}()
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, err := io.ReadAll(reader)
	if err != nil {
		c.logger.Errorln(c.scriptName, err)
		return s
	}

	return string(d)
}

func (c *Cmd) isGBK(data string) bool {
	defer func() {
		err := recover()
		if err != nil {
			c.logger.Errorln(c.scriptName, err)
		}
	}()
	length := len(data)
	var i = 0
	for i < length {
		if data[i] <= 0x7f {
			// 编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			// 大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}
