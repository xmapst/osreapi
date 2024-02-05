//go:build !windows

package exec

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/go-cmd/cmd"
)

func (c *Cmd) selfScriptSuffix() string {
	switch c.shell {
	case "ash":
		return ".ash"
	case "bash":
		return ".bash"
	case "csh":
		return ".csh"
	case "dash":
		return ".dash"
	case "ksh":
		return ".ksh"
	case "shell", "sh":
		return ".sh"
	case "tcsh":
		return ".tcsh"
	case "zsh":
		return ".zsh"
	default:
		return ".bash"
	}
}

func (c *Cmd) beforeExec() []func(cmd *exec.Cmd) {
	return []func(cmd *exec.Cmd){
		func(cmd *exec.Cmd) {
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Setpgid: true,
			}
		},
	}
}

func (c *Cmd) selfCmd() {
	switch c.shell {
	case "ash":
		c.cmd = cmd.NewCmdOptions(c.ops, "ash", "-c", c.scriptName)
	case "csh":
		c.cmd = cmd.NewCmdOptions(c.ops, "csh", "-c", c.scriptName)
	case "dash":
		c.cmd = cmd.NewCmdOptions(c.ops, "dash", "-c", c.scriptName)
	case "ksh":
		c.cmd = cmd.NewCmdOptions(c.ops, "ksh", "-c", c.scriptName)
	case "shell", "sh":
		// 严格模式
		// c.exec = exec.CommandContext(c.ctx, "sh", "-e", c.absFilePath)
		c.cmd = cmd.NewCmdOptions(c.ops, "sh", c.scriptName)
	case "tcsh":
		c.cmd = cmd.NewCmdOptions(c.ops, "tcsh", "-c", c.scriptName)
	case "zsh":
		c.cmd = cmd.NewCmdOptions(c.ops, "zsh", "-c", c.scriptName)
	default:
		// 严格模式
		// c.exec = exec.CommandContext(c.ctx, "bash", "--noprofile", "--norc", "-e", "-o", "pipefail", c.absFilePath)
		// -o pipefail 管道中最后一个返回非零退出状态码的命令的退出状态码将作为该管道命令的返回值，若所有命令的退出状态码都为零则返回零
		c.cmd = cmd.NewCmdOptions(c.ops, "bash", "-o", "pipefail", c.scriptName)
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
		if err = c.cmd.Stop(); err != nil {
			c.logger.Errorln(c.scriptName, err)
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
		// err = errors.New("exec time out")
		// If you use cmd.Process.Kill() directly, only the child process is killed,
		// but the grandchild process is not killed
		// err := cmd.Process.Kill()
		if err = c.cmd.Stop(); err != nil {
			c.logger.Errorln(c.scriptName, err)
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

func (c *Cmd) utf8ToGb2312(s string) string {
	return s
}

func (c *Cmd) transform(line string) string {
	return line
}
