package exec

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/segmentio/ksuid"
)

var (
	ErrTimeOut = errors.New("forced termination by timeout")
	ErrManual  = errors.New("artificial force termination")
)

const (
	// Stop 0, Running 1, Pending 2
	Stop int64 = iota
	Running
	Pending
	Paused

	Killed    = -997
	Timeout   = -998
	SystemErr = -999
)

var StateMap = map[int64]string{
	SystemErr: "System Error",
	Killed:    "Killed",
	Timeout:   "Timeout",
	Stop:      "Stop",
	Running:   "Running",
	Pending:   "Pending",
	Paused:    "Paused",
}

type Cmd struct {
	ops        cmd.Options
	cmd        *cmd.Cmd
	stderrBuf  *cmd.OutputBuffer
	ctx        context.Context
	cancel     context.CancelFunc
	logger     Ilogger
	env        []string
	shell      string
	content    string
	workspace  string
	scriptName string
	scriptDir  string
	timeout    time.Duration
	consoleCh  chan<- string
}

func New(options ...Option) (*Cmd, error) {
	var c = &Cmd{
		ctx:       context.Background(),
		stderrBuf: cmd.NewOutputBuffer(),
		logger:    stdDebugLogger{},
		timeout:   30 * time.Minute,
		workspace: filepath.Join(os.TempDir(), "workspace"),
		scriptDir: filepath.Join(os.TempDir(), "script"),
	}
	for _, option := range options {
		option(c)
	}
	c.ops = cmd.Options{
		Buffered:       false,
		Streaming:      true,
		BeforeExec:     c.beforeExec(),
		LineBufferSize: 491520,
	}
	c.scriptName = filepath.Join(c.scriptDir, ksuid.New().String())
	c.scriptName = c.scriptName + c.scriptSuffix()
	if err := os.MkdirAll(c.scriptDir, os.ModePerm); err != nil {
		c.logger.Errorln(c.scriptName, err)
		return nil, err
	}
	c.logger.Infoln(c.scriptName, "create script")
	if c.shell == "cmd" || c.shell == "powershell" {
		c.content = c.utf8ToGb2312(c.content)
	}
	if err := os.WriteFile(c.scriptName, []byte(c.content), os.ModePerm); err != nil {
		c.logger.Errorln(c.scriptName, err)
		return nil, err
	}
	return c, nil
}

func (c *Cmd) scriptSuffix() string {
	switch c.shell {
	case "python", "python2", "python3", "py", "py2", "py3":
		return ".py"
	}
	return c.selfScriptSuffix()
}

func (c *Cmd) clear() {
	// clear tmp script
	c.logger.Infoln(c.scriptName, "cleanup script")
	if err := os.Remove(c.scriptName); err != nil {
		c.logger.Errorln(c.scriptName, err)
	}
	c.logger.Infoln(c.scriptName, "cleanup script dir")
	if err := os.RemoveAll(c.scriptDir); err != nil {
		c.logger.Errorln(c.scriptName, err)
	}
}

func (c *Cmd) newCmd() {
	c.ctx, c.cancel = context.WithTimeout(context.Background(), c.timeout)
	switch c.shell {
	case "python", "python2", "py2", "py":
		c.cmd = cmd.NewCmdOptions(c.ops, "python2", c.scriptName)
	case "python3", "py3":
		c.cmd = cmd.NewCmdOptions(c.ops, "python3", c.scriptName)
	default:
		c.selfCmd()
	}

	// inject env
	c.logger.Infoln(c.scriptName, "inject customize env", c.env)
	c.cmd.Env = append(append(os.Environ(), c.env...),
		fmt.Sprintf("TASK_WORKSPACE=%s", c.workspace),
		fmt.Sprintf("TASK_STEP_SCRIPT_DIR=%s", c.scriptDir),
		fmt.Sprintf("TASK_STEP_SCRIPT_PATH=%s", c.scriptName),
	)

	// set workspace
	c.cmd.Dir = c.workspace
}

func (c *Cmd) consoleOutput() {
	defer func() {
		close(c.consoleCh)
		c.logger.Debugln("stop console print")
	}()
	for {
		var line string
		var open bool
		select {
		case <-c.ctx.Done():
			return
		case line, open = <-c.cmd.Stdout:
			if !open {
				c.cmd.Stdout = nil
				continue
			}
		case line, open = <-c.cmd.Stderr:
			if !open {
				c.cmd.Stderr = nil
				continue
			}
			_, _ = c.stderrBuf.Write([]byte(line))
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = c.transform(line)

		c.logger.Infoln(c.scriptName, line)
		if c.consoleCh != nil {
			c.consoleCh <- line
		}
	}
}

func last(slice []string) string {
	if len(slice) > 0 {
		return slice[len(slice)-1]
	}
	return ""
}
