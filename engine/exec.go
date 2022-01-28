package engine

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/osreapi/utils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type Cmd struct {
	Shell           string
	Name            string
	TaskID          string
	Step            int
	Content         string
	ExternalEnvVars []string
	ExecTimeout     time.Duration
	absFilePath     string
	exec            *exec.Cmd
	context         context.Context
	cancelFunc      context.CancelFunc
	Log             *logrus.Entry
}

func (c *Cmd) Initial() error {
	c.absFilePath = filepath.Join(utils.TmpDir, c.Name)
	switch c.Shell {
	case "sh", "shell", "bash", "ash":
		c.absFilePath = c.absFilePath + ".sh"
	case "cmd", "bat":
		c.absFilePath = c.absFilePath + ".bat"
	case "powershell", "ps", "ps1":
		c.absFilePath = c.absFilePath + ".ps1"
	case "python", "python2", "python3", "py", "py2", "py3":
		c.absFilePath = c.absFilePath + ".py"
	}
	c.Log.Infof("create script %s", filepath.Base(c.absFilePath))
	return os.WriteFile(c.absFilePath, []byte(c.Content), 0777)
}

func (c *Cmd) clear() {
	// clear tmp script

	c.Log.Infof("cleanup script %s", filepath.Base(c.absFilePath))
	_ = os.Remove(c.absFilePath)
}

func (c *Cmd) initCmd(execID string) bool {
	c.Log.Infof("set exec timeout %s", c.ExecTimeout.String())
	c.context, c.cancelFunc = context.WithTimeout(context.Background(), c.ExecTimeout)
	switch c.Shell {
	case "shell", "sh":
		if runtime.GOOS == "windows" {
			return false
		}
		c.exec = exec.CommandContext(c.context, "sh", "-c", c.absFilePath)
	case "ash":
		if runtime.GOOS == "windows" {
			return false
		}
		c.exec = exec.CommandContext(c.context, "ash", "-c", c.absFilePath)
	case "bash":
		if runtime.GOOS == "windows" {
			return false
		}
		c.exec = exec.CommandContext(c.context, "bash", "-c", c.absFilePath)
	case "python", "python2", "py2", "py":
		c.exec = exec.CommandContext(c.context, "python2", c.absFilePath)
	case "python3", "py3":
		c.exec = exec.CommandContext(c.context, "python3", c.absFilePath)
	case "cmd", "bat":
		if runtime.GOOS != "windows" {
			return false
		}
		c.exec = exec.CommandContext(c.context, "cmd", "/C", c.absFilePath)
	case "powershell", "ps", "ps1":
		if runtime.GOOS != "windows" {
			return false
		}
		c.exec = exec.CommandContext(c.context, "powershell", "-NoLogo", "-NonInteractive", "-File", c.absFilePath)
	default:
		return false
	}
	c.exec.Env = append(append(os.Environ(), c.ExternalEnvVars...), fmt.Sprintf("TASK_ID=%s", execID))
	return true
}
